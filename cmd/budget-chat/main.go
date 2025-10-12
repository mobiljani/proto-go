package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	port := 8080

	list, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	fmt.Printf("Server started on port %d\n", port)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	for {
		conn, err := list.Accept()

		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		go handleConnection(conn)
	}
}

func handleConnection(connection net.Conn) {
	var name string

	defer connection.Close()

	fmt.Printf("Client connected: %s\n", connection.RemoteAddr().String())
	connection.Write([]byte("Welcome to budgetchat! What shall I call you?\n"))
	scanner := bufio.NewScanner(connection)

	for scanner.Scan() {
		rb := scanner.Bytes()
		in := string(rb)
		fmt.Printf("Received: %s \n", rb)

		if name == "" {
			if !validateName(in) {
				connection.Write([]byte("Name must be between 1 and 16 characters\n"))
				break
			}
			name = in
			continue
		}

		m := fmt.Sprintf("[%s] %s\n", name, in)
		connection.Write([]byte(m))
	}
}

func validateName(name string) bool {
	return len(name) >= 1 && len(name) <= 16
}
