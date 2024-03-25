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

	start_line := lines[0]
	parts := strings.Split(start_line, " ")
	if len(parts) < 3 {
		fmt.Println("Startline malformed!")
		os.Exit(1)
	}

	path := parts[1]
	route(path, c)
}

func route(path string, c net.Conn) {
	if path == "/" {
		msg := "HTTP/1.1 200 OK\r\n\r\n"
		c.Write([]byte(msg))
		os.Exit(0)
	}

	path_parts := strings.SplitN(path, "/", 3)
	if len(path_parts) > 1 && path_parts[1] == "echo" {
		content := path_parts[2]
		content_length := len(content)
		msg := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", content_length, content)
		fmt.Println(msg)
		c.Write([]byte(msg))
		os.Exit(0)
	}

	msg := "HTTP/1.1 404 NOT FOUND\r\n\r\n"
	c.Write([]byte(msg))
	os.Exit(1)
}
