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
		in := scanner.Bytes()

		// https://pkg.go.dev/fmt

		/*

			Message: [121] - [0 0 60 0] : [0 0 100 0] String Q0@ Human Q 12288 16384
			Message: [111] - [0 0 240 0] : [0 0 0 5] String I� Human I 40960 5
			Message: [111] - [0 0 60 73] : [0 0 0 144] String I0;d Human I 12347 100
			Message: [111] - [0 0 60 72] : [0 0 0 146] String I0:f Human I 12346 102
			Message: [111] - [0 0 60 71] : [0 0 0 145] String I09e Human I 12345 101

		*/

		fmt.Printf("Message: %v - %v : %v  String %s Human %s %d %d \n", in[0:1], in[1:5], in[5:9], string(in), string(in[0:1]), binary.BigEndian.Uint32(in[1:5]), binary.BigEndian.Uint32(in[5:9]))

		if string(in[0]) == "I" {
			// TODO: Price can be negative
			new := entry{time: binary.BigEndian.Uint32(in[1:5]), price: binary.BigEndian.Uint32(in[5:9])}
			list = append(list, new)
			fmt.Print(list)
		}

		if string(in[0]) == "Q" {
			from := binary.BigEndian.Uint32(in[1:5])
			to := binary.BigEndian.Uint32(in[5:9])
			var count, total, mean int

			for _, item := range list {
				if item.time >= from && item.time <= to {
					count = count + 1
					total = total + int(item.price)
				}
			}

			if count > 0 {
				mean = total / count
			}

			out := make([]byte, 4)

			binary.BigEndian.PutUint32(out, uint32(mean))
			connection.Write(out)

			fmt.Printf("Mean price is %d - %v\n", mean, out)
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
