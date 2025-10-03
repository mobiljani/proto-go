package main

import (
	"fmt"
	"net"
	"time"
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
	fmt.Printf("Connection from %s\n", connection.RemoteAddr().String())
	connection.SetReadDeadline(time.Now().Add(time.Minute))

	buffer := make([]byte, 1024)

	length, err := connection.Read(buffer)

	if err != nil {
		fmt.Printf("Failed to read message from client: %v\n", err)
		connection.Close()
		return
	}

	message := string(buffer[:length])

	fmt.Printf("Received: %d bytes \t: %s\n", length, message)

	connection.Write([]byte(message))

	//connection.Close()

}
