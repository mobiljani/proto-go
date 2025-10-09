package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"net"
	"os"
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
	port := 8081

	// Configure slog to output to stdout with text format
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	list, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	slog.Info("Server started", "port", port)

	if err != nil {
		slog.Info(err.Error(), "error", err)
		return
	}
	for {
		conn, err := list.Accept()

		if err != nil {
			slog.Info(err.Error(), "error", err)
			return
		}

		go handleConnection(conn)
	}
}

func handleConnection(connection net.Conn) {
	defer connection.Close()

	slog.Info("Client connected", "address", connection.RemoteAddr().String())

	buffer := make([]byte, 4*1024)
	for {
		bytesRead, err := connection.Read(buffer)
		if err != nil {
			slog.Info("Error during read", "error", err)
			break
		}

		slog.Info("Message", "msg", string(buffer[0:bytesRead]))

		req := request{}

		jsonError := json.Unmarshal(buffer[0:bytesRead], &req)

		if jsonError != nil {
			slog.Info("Error during JSON unmarshal", "error", jsonError)
			break
		}

		if req.Method != "isPrime" {
			slog.Info("Unknown method", "method", req.Method)
			break
		}

		res := response{}

		res.Method = "isPrime"
		res.Prime = isPrime(req.Number)

		s, _ := json.Marshal(res)

		connection.Write([]byte(string(s)))
	}

}

func isPrime(num int) bool {
	nig := big.NewInt(int64(num))
	return nig.ProbablyPrime(0)
}
