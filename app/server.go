package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	c, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	buf := make([]byte, 1024)
	_, err = c.Read(buf)
	fmt.Println("Read incoming request:\n", string(buf))

	lines := strings.Split(string(buf), "\r\n")
	if len(lines) < 1 {
		fmt.Println("Startline missing!")
		os.Exit(1)
	}

	startLine := lines[0]
	parts := strings.Split(startLine, " ")
	if len(parts) < 3 {
		fmt.Println("Startline malformed!")
		os.Exit(1)
	}

	headers := parseHeaders(lines[1:])
	path := parts[1]
	route(c, path, headers)
}

func parseHeaders(lines []string) map[string]string {
	headers := make(map[string]string)
	for _, l := range lines {
		parsed := strings.SplitN(l, ": ", 2)
		headers[parsed[0]] = parsed[1]
	}

	return headers
}

func route(c net.Conn, path string, headers map[string]string) {
	if path == "/" {
		msg := "HTTP/1.1 200 OK\r\n\r\n"
		c.Write([]byte(msg))
		os.Exit(0)
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

	handleNotFound(c)
}

func handleEcho(c net.Conn, pathParts []string) {
	content := pathParts[2]
	contentLength := len(content)
	msg := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", contentLength, content)
	fmt.Println(msg)
	c.Write([]byte(msg))
	os.Exit(0)
}

func handleUserAgent(c net.Conn, headers map[string]string) {
	fmt.Sprintf("Not implemented yet!")
	os.Exit(1)
}

func handleNotFound(c net.Conn) {
	msg := "HTTP/1.1 404 NOT FOUND\r\n\r\n"
	c.Write([]byte(msg))
	os.Exit(1)
}
