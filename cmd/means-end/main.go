package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"math/big"
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
		time  int
		price int
	}

	var list []entry

	fmt.Printf("Client connected: %s\n", connection.RemoteAddr().String())

	scanner := bufio.NewScanner(connection)
	scanner.Split(everyNineBytes)
	for scanner.Scan() {
		rb := scanner.Bytes()

		first := new(big.Int)
		second := new(big.Int)

		first.SetBytes(rb[1:5])
		second.SetBytes(rb[5:9])

		//fmt.Printf("Message: %v - %v : %v  String %s Human %s %d %d \n", rb[0:1], rb[1:5], rb[5:9], string(rb), string(rb[0:1]), binary.BigEndian.Uint32(rb[1:5]), binary.BigEndian.Uint32(rb[5:9]))

		if string(rb[0]) == "I" {
			new := entry{time: int(first.Int64()), price: int(second.Int64())}
			list = append(list, new)
		}

		if string(rb[0]) == "Q" {
			//[4intoverflow.test] FAIL:Q -2100000000 210000000: expected 2050000000 (6150000000/3), got 0

			var count, mean int
			var total int64

			fmt.Printf("Query from %d to %d ", int(first.Int64()), int(second.Int64()))

			for _, item := range list {
				if item.time >= int(first.Int64()) && item.time <= int(second.Int64()) {
					count = count + 1
					total += int64(item.price)
				}
			}

			if count > 0 {
				mean = int(total / int64(count))
			}

			wb := make([]byte, 4)

			binary.BigEndian.PutUint32(wb, uint32(mean))
			connection.Write(wb)

			fmt.Printf("%v -> %v \n", rb, wb)

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
