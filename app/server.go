package main

import (
	"bytes"
	"flag"
	"fmt"
	"net"
	"os"
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
	defer c.Close()

	buf := make([]byte, 1024)
	_, err := c.Read(buf)
	if err != nil {
		fmt.Println("Error reading from connection: ", err.Error())
		respondBadRequest(c)
		return
	}
	buf = bytes.Trim(buf, "\x00")

	request, err := parseRequest(buf)
	if err != nil {
		fmt.Println("Error parsing request: ", err.Error())
		respondBadRequest(c)
		return
	}

	route(c, request)
}
