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
	defer connection.Close()

	fmt.Printf("Client connected: %s\n", connection.RemoteAddr().String())

	var name string
	var nameAsked bool

	connection.Write([]byte("Welcome to budgetchat! What shall I call you?"))
	nameAsked = true

	scanner := bufio.NewScanner(connection)
	for scanner.Scan() {
		rb := scanner.Bytes()
		in := string(rb)
		fmt.Printf("Received: %s \n", rb)

		if nameAsked && name == "" {
			name = in
			continue
		}

		m := fmt.Sprintf("[%s] %s\n", name, in)

		connection.Write([]byte(m))

	}

}
