package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
)

const directoryFlag = "directory"
const portFlag = "port"

const defaultPort = 4221

var directory *string
var port *int

const maxRequestSize = 1024

func main() {
	parseArgs()

	l := listen()
	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConnection(c)
	}
}

func parseArgs() {
	directory = flag.String(directoryFlag, "", "Directory to take files from")
	port = flag.Int(portFlag, defaultPort, "Port to listen on")
	flag.Parse()
}

func listen() net.Listener {
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *port))
	if err != nil {
		fmt.Println("Failed to bind to port ", *port)
		os.Exit(1)
	}

	return l
}

func handleConnection(c net.Conn) {
	defer c.Close()

	buf, err := readRequest(c)
	if err != nil {
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

func readRequest(c net.Conn) (buf []byte, err error) {
	buf = make([]byte, maxRequestSize)
	_, e := c.Read(buf)
	if e != nil {
		fmt.Println("Error reading from connection: ", err.Error())
		err = errors.New("error reading from connection")
	}
	buf = bytes.Trim(buf, "\x00")
	return
}
