package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	resolver, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		fmt.Println("Failed to resolve UDP connection")
		os.Exit(1)
	}
	conn, err := net.DialUDP("udp", nil, resolver)
	if err != nil {
		fmt.Println("Failed to dial UDP connection")
		os.Exit(1)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println(">")
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading message: %v\n", err)
			os.Exit(1)
		}
		_, err = conn.Write([]byte(msg))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error sending message: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("sent: %s", msg)
	}
}
