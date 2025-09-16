package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

const BUFFER_SIZE = 8

type RequestState int

const (
	stateInitialised RequestState = iota
	stateDone
)

type Request struct {
	RequestLine RequestLine
	state       RequestState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, BUFFER_SIZE, BUFFER_SIZE)
	readIndex := 0

	request := &Request{
		state: stateInitialised,
	}

	for request.state != stateDone {
		if readIndex >= len(buf) {
			biggerBuff := make([]byte, len(buf)*2)
			copy(biggerBuff, buf)
			buf = biggerBuff
		}
		bytesRead, err := reader.Read(buf[readIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				request.state = stateDone
			}
			return nil, err
		}
		readIndex += bytesRead

		bytesParsed, err := request.parse(buf[:readIndex])
		if err != nil {
			return nil, err
		}
		copy(buf, buf[bytesParsed:])
		readIndex -= bytesParsed
	}
	return request, nil
}

func parseRequestLine(d []byte) (*RequestLine, int, error) {
	index := bytes.Index(d, []byte("\r\n"))
	if index == -1 {
		return nil, 0, nil
	}
	reqLine := string(d[:index])
	requestLine, err := requestLineFromString(reqLine)
	if err != nil {
		return nil, 0, err
	}

	return requestLine, index + 2, nil

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

func (r *Request) parse(data []byte) (int, error) {
	switch r.state {
	case stateInitialised:
		requestLine, n, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, nil
		}
		r.RequestLine = *requestLine
		r.state = stateDone
		return n, nil
	case stateDone:
		return 0, fmt.Errorf("error: trying to read data in a done state")
	default:
		return 0, fmt.Errorf("unknown state")
	}
}
