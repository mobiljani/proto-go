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

	type entry struct {
		time  uint32
		price uint32
	}

	var list []entry

	fmt.Printf("Client connected: %s\n", connection.RemoteAddr().String())

	scanner := bufio.NewScanner(connection)
	scanner.Split(everyNineBytes)
	for scanner.Scan() {
		bytes := scanner.Bytes()

		// https://pkg.go.dev/fmt

		/*

			Message: [121] - [0 0 60 0] : [0 0 100 0] String Q0@ Human Q 12288 16384
			Message: [111] - [0 0 240 0] : [0 0 0 5] String I� Human I 40960 5
			Message: [111] - [0 0 60 73] : [0 0 0 144] String I0;d Human I 12347 100
			Message: [111] - [0 0 60 72] : [0 0 0 146] String I0:f Human I 12346 102
			Message: [111] - [0 0 60 71] : [0 0 0 145] String I09e Human I 12345 101

		*/

		fmt.Printf("Message: %o - %o : %o  String %s Human %s %d %d \n", bytes[0:1], bytes[1:5], bytes[5:9], string(bytes), string(bytes[0:1]), binary.BigEndian.Uint32(bytes[1:5]), binary.BigEndian.Uint32(bytes[5:9]))

		if string(bytes[0:1]) == "I" {
			// TODO: Price can be negative
			new := entry{time: binary.BigEndian.Uint32(bytes[1:5]), price: binary.BigEndian.Uint32(bytes[5:9])}
			list = append(list, new)
			fmt.Print(list)
		}

		if string(bytes[0:1]) == "Q" {
			connection.Write([]byte("todo"))
		}

	}

}

func everyNineBytes(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF {
		return 0, nil, bufio.ErrFinalToken
	}

	if len(data) < 9 {
		return 0, nil, nil
	}

	return 9, data[:9], nil
}
