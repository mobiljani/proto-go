package main

import (
	"bufio"
	"encoding/binary"
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

		fmt.Printf("Message: %o - %o : %o  String %s Human %s %d %d \n", bytes[0:1], bytes[1:5], bytes[5:9], bytes, string(bytes[0:1]), binary.BigEndian.Uint32(bytes[1:5]), binary.BigEndian.Uint32(bytes[5:9]))

	}

}

func everyNineBytes(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if len(data) < 9 || atEOF {
		return 0, nil, bufio.ErrFinalToken
	}

	return 9, data[:9], nil
}
