package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"path"
	"strings"
)

const USER_AGENT_HEADER = "User-Agent"

const DIRECTORY_FLAG = "directory"
const PORT_FLAG = "port"

const EMPTY_DIR = "/var/empty/"
const DEFAULT_PORT = 4221

var directory *string
var port *int

func main() {
	directory = flag.String(DIRECTORY_FLAG, "", "Directory to take files from")
	port = flag.Int(PORT_FLAG, DEFAULT_PORT, "Port to listen on")
	flag.Parse()

	l := listen(*port)
	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConnection(c)
	}
}

func listen(port int) net.Listener {
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		fmt.Println("Failed to bind to port ", port)
		os.Exit(1)
	}

	return l
}

func handleConnection(c net.Conn) {
	buf := make([]byte, 1024)
	_, err := c.Read(buf)
	if err != nil {
		fmt.Println("Error reading from connection: ", err.Error())
		respondBadRequest(c)
		return
	}

	request, err := parseRequest(buf)
	if err != nil {
		fmt.Println("Error parsing request: ", err.Error())
		respondBadRequest(c)
		return
	}

	route(c, request)
}

const rootRoute = "/"
const echoRoute = "/echo/"
const userAgentRoute = "/user-agent"
const filesPath = "/files/"

const get = "GET"
const post = "POST"

func route(c net.Conn, request Request) {
	switch {
	case exactMatch(request, rootRoute, get):
		handleRoot(c)
	case match(request, echoRoute, get):
		handleEcho(c, request)
	case match(request, userAgentRoute, get):
		handleUserAgent(c, request)
	case match(request, filesPath, get):
		handleGetFiles(c, request)
	case match(request, filesPath, post):
		handlePostFiles(c, request)

	default:
		respondNotFound(c)
	}
}

func handleRoot(c net.Conn) {
	msg := "HTTP/1.1 200 OK\r\n\r\n"
	c.Write([]byte(msg))
	c.Close()
}

func handleEcho(c net.Conn, request Request) {
	content, _ := strings.CutPrefix(request.uri, echoRoute)
	contentLength := len(content)

	msg := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", contentLength, content)
	c.Write([]byte(msg))
	c.Close()
}

func handleUserAgent(c net.Conn, request Request) {
	userAgent, ok := request.headers[USER_AGENT_HEADER]

	if !ok {
		fmt.Sprintf("No User Agent header provided!")
		respondBadRequest(c)
		return
	}

	contentLength := len(userAgent)
	msg := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", contentLength, userAgent)
	c.Write([]byte(msg))
	c.Close()
}

func handleGetFiles(c net.Conn, request Request) {
	fileName, valid := strings.CutPrefix(request.uri, filesPath)
	if !valid {
		fmt.Println("No file name given!")
		respondBadRequest(c)
		return
	}

	filePath := path.Join(*directory, fileName)
	_, err := os.Stat(filePath)

	if err != nil {
		respondNotFound(c)
		return
	}

	dat, err := os.ReadFile(filePath)
	if err != nil {
		respondNotFound(c)
		return
	}

	msg := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n", len(dat))
	c.Write([]byte(msg))
	c.Write(dat)
	c.Close()
}

func handlePostFiles(c net.Conn, request Request) {
	fileName, valid := strings.CutPrefix(request.uri, filesPath)
	if !valid {
		fmt.Println("No file name given!")
		respondBadRequest(c)
		return
	}

	filePath := path.Join(*directory, fileName)

	os.WriteFile(filePath, []byte(request.body), 0644)
	msg := "HTTP/1.1 201 ACCEPTED\r\n\r\n"
	c.Write([]byte(msg))
	c.Close()
}

func respondNotFound(c net.Conn) {
	msg := "HTTP/1.1 404 NOT FOUND\r\n\r\n"
	c.Write([]byte(msg))
	c.Close()
}

func respondBadRequest(c net.Conn) {
	msg := "HTTP/1.1 400 BAD REQUEST\r\n\r\n"
	c.Write([]byte(msg))
	c.Close()
}

func exactMatch(request Request, uri string, method string) bool {
	return request.uri == uri && request.method == method
}

func match(request Request, prefix string, method string) bool {
	return strings.HasPrefix(request.uri, prefix) && request.method == method
}
