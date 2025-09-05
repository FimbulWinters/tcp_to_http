package main

import (
	"fmt"
	"io"
	"net"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Println(err)
	}
	defer listener.Close()

	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("connection accepted")
		fileChan := getLinesChannel(connection)
		for line := range fileChan {
			fmt.Println("read:", line)
		}

	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	fileChan := make(chan string)
	go func() {
		defer f.Close()
		defer close(fileChan)
		fileData := make([]byte, 8)
		currentLine := ""
		for {
			data, err := f.Read(fileData)
			if data > 0 {
				currentLine += string(fileData[:data])
				split := strings.Split(currentLine, "\n")
				if len(split) > 1 {
					part := strings.Join(split[:len(split)-1], " ")
					fileChan <- part
					currentLine = split[len(split)-1]
				}
			}
			if err != nil {
				if err == io.EOF {
					break
				}
			}
		}
	}()

	return fileChan
}
