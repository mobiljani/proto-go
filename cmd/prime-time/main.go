package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/big"
	"net"
)

type request struct {
	Method string `json:"method"`
	Number int    `json:"number"`
}

type response struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

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
	for scanner.Scan() {
		bytes := scanner.Bytes()

		fmt.Printf("Message: %s\n", string(bytes))

		req := request{}

		jsonError := json.Unmarshal(bytes, &req)

		if jsonError != nil {
			fmt.Printf("Error during JSON unmarshal: %v\n", jsonError)
			connection.Write([]byte(string("meh")))
			break
		}

		if req.Method != "isPrime" {
			fmt.Printf("Unknown method: %s\n", req.Method)
			connection.Write([]byte(string("meh")))
			break
		}

		res := response{}

		res.Method = "isPrime"
		res.Prime = isPrime(req.Number)

		s, _ := json.Marshal(res)

		connection.Write([]byte(string(s) + "\n"))
		fmt.Printf("Response: %s\n", string(s))
	}

}

func isPrime(num int) bool {
	nig := big.NewInt(int64(num))
	return nig.ProbablyPrime(0)
}
