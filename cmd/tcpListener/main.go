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
	lines := make(chan string)

	go func() {
		defer f.Close()
		defer close(lines)

		buf := make([]byte, 8)
		var pending string

		for {
			n, err := f.Read(buf)
			if n > 0 {
				chunk := string(buf[:n])
				parts := strings.Split(chunk, "\n")

				for i := 0; i < len(parts)-1; i++ {
					lines <- pending + parts[i]
					pending = ""
				}
				pending += parts[len(parts)-1]
			}

			if err != nil {
				if pending != "" {
					lines <- pending
				}

				break
			}
		}
	}()

	return lines
}
