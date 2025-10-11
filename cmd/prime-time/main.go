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
	Number *int   `json:"number"`
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

		if jsonError != nil || req.Method != "isPrime" || req.Number == nil {
			connection.Write([]byte("meh"))
			break
		}

		s, _ := json.Marshal(response{Method: "isPrime", Prime: isPrime(*req.Number)})
		connection.Write([]byte(string(s) + "\n"))
		fmt.Printf("Response: %s\n", string(s))
	}

}

func isPrime(num int) bool {
	nig := big.NewInt(int64(num))
	return nig.ProbablyPrime(0)
}
