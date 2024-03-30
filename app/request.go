package main

import (
	"errors"
	"strings"
)

// see https://datatracker.ietf.org/doc/html/rfc2616#section-4

type Request struct {
	method  string
	uri     string
	headers map[string]string
	body    string
}

type requestLine struct {
	method  string
	uri     string
	version string
}

func parseRequest(buf []byte) (req Request, err error) {
	head, body, valid := strings.Cut(string(buf), "\r\n\r\n")

	if !valid {
		err = errors.New("HTTP Request Malformed (must have a double CRLF cut)")
		return
	}

	reqLine, headerLines, _ := strings.Cut(head, "\r\n")
	parsedRequestLine, _ := parseRequestLine(reqLine)
	headers := parseHeaders(strings.Split(headerLines, "\r\n"))

	req = Request{
		method:  parsedRequestLine.method,
		uri:     parsedRequestLine.uri,
		headers: headers,
		body:    body,
	}
	return
}

func parseRequestLine(line string) (parsedRequestLine requestLine, err error) {
	parts := strings.Split(line, " ")

	if len(parts) != 3 {
		err = errors.New("HTTP Request Line Malformed (must have exactly 3 tokens)")
		return
	}

	parsedRequestLine = requestLine{parts[0], parts[1], parts[2]}
	return
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
