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

	scanner := bufio.NewScanner(connection)
	scanner.Split(everyNineBytes)
	for scanner.Scan() {
		bytes := scanner.Bytes()

		// https://pkg.go.dev/fmt

		fmt.Printf("%s Message: %o - %o : %o  String %s \n", connection.RemoteAddr().String(), bytes[0:1], bytes[1:5], bytes[5:9], bytes)

	}

}

func everyNineBytes(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if len(data) < 9 || atEOF {
		return 0, nil, nil
	}

	return 9, data[:9], nil
}
