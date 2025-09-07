package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	d, err := io.ReadAll(reader)
	if err != nil {
		fmt.Println("Failed to read data from provider reader")
		return nil, fmt.Errorf("failed to read data from provider reader: %w", err)
	}
	requestLine, err := parseRequestLine(d)
	if err != nil {
		return nil, fmt.Errorf("could not parse request line: %w", err)
	}

	request := &Request{
		RequestLine: *requestLine,
	}
	return request, nil
}

func parseRequestLine(d []byte) (*RequestLine, error) {
	index := bytes.Index(d, []byte("\r\n"))
	if index == -1 {
		return nil, fmt.Errorf("Couldn't find CLRF char")
	}
	reqLine := string(d[:index])
	requestLine, err := requestLineFromString(reqLine)
	if err != nil {
		return nil, err
	}

	return requestLine, nil

}

func requestLineFromString(reqLine string) (*RequestLine, error) {
	parts := strings.Split(reqLine, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("Request line missing mandatory part")
	}
	method := parts[0]
	for _, c := range method {
		if c < 'A' || c > 'Z' {
			return nil, fmt.Errorf("invalid method: %s", method)
		}
	}
	requestTarget := parts[1]
	httpVersionparts := strings.Split(parts[2], "/")
	if len(httpVersionparts) != 2 {
		return nil, fmt.Errorf("incorrectly formatted http element")
	}
	if httpVersionparts[0] != "HTTP" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", httpVersionparts[0])
	}
	version := httpVersionparts[1]
	if version != "1.1" {
		return nil, fmt.Errorf("wrong http version number")
	}
	return &RequestLine{
		Method:        method,
		RequestTarget: requestTarget,
		HttpVersion:   version,
	}, nil
}
