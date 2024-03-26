package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

const USER_AGENT_HEADER = "User-Agent"

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	// Probably want something like...
	// Infinite loop to accept connections and spawn worker threads
	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConnection(c)
	}
}

func handleConnection(c net.Conn) {
	buf := make([]byte, 1024)
	_, err := c.Read(buf)
	if err != nil {
		fmt.Println("Error reading from connection: ", err.Error())
		return
	}

	fmt.Println("Read incoming request:\n", string(buf))

	lines := strings.Split(string(buf), "\r\n")
	if len(lines) < 1 {
		fmt.Println("Startline missing!")
		writeBadRequest(c)
		return
	}

	startLine := lines[0]
	parts := strings.Split(startLine, " ")
	if len(parts) < 3 {
		fmt.Println("Startline malformed!")
		writeBadRequest(c)
		return
	}

	headers := parseHeaders(lines[1:])
	path := parts[1]
	route(c, path, headers)

}

func parseHeaders(lines []string) map[string]string {
	headers := make(map[string]string)
	for _, l := range lines {
		if len(l) == 0 {
			break
		}

		parsed := strings.SplitN(l, ": ", 2)
		headers[parsed[0]] = parsed[1]
	}

	return headers
}

func route(c net.Conn, path string, headers map[string]string) {
	if path == "/" {
		msg := "HTTP/1.1 200 OK\r\n\r\n"
		c.Write([]byte(msg))
		c.Close()
	}

	pathParts := strings.SplitN(path, "/", 3)
	if len(pathParts) > 1 {
		switch pathParts[1] {
		case "echo":
			handleEcho(c, pathParts)
		case "user-agent":
			handleUserAgent(c, headers)
		}
	}

	writeNotFound(c)
	c.Close()
}

func handleEcho(c net.Conn, pathParts []string) {
	content := pathParts[2]
	contentLength := len(content)
	msg := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", contentLength, content)
	fmt.Println(msg)
	c.Write([]byte(msg))
	c.Close()
}

func handleUserAgent(c net.Conn, headers map[string]string) {
	userAgent, ok := headers[USER_AGENT_HEADER]

	if !ok {
		fmt.Sprintf("No User Agent header provided!")
		writeBadRequest(c)
		c.Close()
		return
	}

	contentLength := len(userAgent)
	msg := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", contentLength, userAgent)
	fmt.Println(msg)
	c.Write([]byte(msg))
	c.Close()
}

func writeNotFound(c net.Conn) {
	msg := "HTTP/1.1 404 NOT FOUND\r\n\r\n"
	c.Write([]byte(msg))
}

func writeBadRequest(c net.Conn) {
	msg := "HTTP/1.1 400 BAD REQUEST\r\n\r\n"
	c.Write([]byte(msg))
}
