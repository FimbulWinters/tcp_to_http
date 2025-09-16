package main

import (
	"fmt"
	"net"

	"github.com/FimbulWinters/tcp_to_http/internal/request"
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

		fileChan, err := request.RequestFromReader(connection)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("Request line: \n - Method: %s\n - Target: %s\n - Version: %s", fileChan.RequestLine.Method, fileChan.RequestLine.RequestTarget, fileChan.RequestLine.HttpVersion)

	}
}
