package main

import (
	"fmt"
	"io"
	"net"
)

func main() {
	port := 8080
	fmt.Printf("Starting server on port %d\n", port)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		fmt.Printf("Failed to create a server: %v\n", err)
		return
	}

	for {
		connection, err := listener.Accept()

		if err != nil {
			fmt.Printf("Connection error %v\n", err)
			continue
		}

		go handleConnection(connection)
	}

}

func handleConnection(connection net.Conn) {
	defer connection.Close()

	// _, err := io.Copy(connection, connection)
	// if err != nil {
	// 	fmt.Printf("Error during io.Copy: %v\n", err)
	// }

	fmt.Printf("Connection from %s\n", connection.RemoteAddr().String())

	buffer := make([]byte, 32*1024)

	for {

		length, err := connection.Read(buffer)

		if err != nil {
			fmt.Printf("Failed to read message from client: %v\n", err)
			if err == io.EOF {
				fmt.Printf("Client closed the connection\n")
				connection.Write(buffer)
			}
			connection.Close()
			return
		}

		if length == 0 {
			fmt.Printf("Empty message\n")
			continue
		}

		message := string(buffer[:length])
		fmt.Printf("Received: %s\n", message)

	}

}
