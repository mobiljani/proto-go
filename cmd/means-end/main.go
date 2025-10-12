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
		rb := scanner.Bytes()

		//fmt.Printf("Message: %v - %v : %v  String %s Human %s %d %d \n", rb[0:1], rb[1:5], rb[5:9], string(rb), string(rb[0:1]), binary.BigEndian.Uint32(rb[1:5]), binary.BigEndian.Uint32(rb[5:9]))

		if string(rb[0]) == "I" {
			// TODO: Price can be negative
			new := entry{time: binary.BigEndian.Uint32(rb[1:5]), price: binary.BigEndian.Uint32(rb[5:9])}
			list = append(list, new)
			//fmt.Print(list)
		}

		if string(rb[0]) == "Q" {
			from := binary.BigEndian.Uint32(rb[1:5])
			to := binary.BigEndian.Uint32(rb[5:9])
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

			wb := make([]byte, 4)

			binary.BigEndian.PutUint32(wb, uint32(mean))
			connection.Write(wb)

			fmt.Printf("Mean price is %d - %v\n", mean, wb)
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
