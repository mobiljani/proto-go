package main

import (
	"bufio"
	"context"
	"fmt"
	"net"

	"github.com/maniartech/signals"
)

var userMessaged = signals.New[string]()
var serverMessaged = signals.New[string]()

func main() {

	list, err := net.Listen("tcp", ":8080")
	fmt.Printf("Server started on port 8080\n")

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

		go handleUserConnection(conn)

	}
}

func handleUserConnection(connection net.Conn) {
	defer connection.Close()

	serverMessaged.AddListener(func(ctx context.Context, msg string) {
		fmt.Printf("Downstream msg to user: '%s'\n", msg)
		_, err := connection.Write([]byte(msg))
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
	})

	downstream, err := net.Dial("tcp", "chat.protohackers.com:16963")
	fmt.Printf("Client started on port chat.protohackers.com:16963\n")

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	scanner := bufio.NewScanner(connection)
	fmt.Printf("Starting to read user messages\n")
	go handleServerConnection(downstream)

	for scanner.Scan() {
		in := scanner.Text()
		fmt.Printf("user: '%s'\n", in)
		userMessaged.Emit(context.Background(), in)
	}
}

func handleServerConnection(downstream net.Conn) {
	defer downstream.Close()

	userMessaged.AddListener(func(ctx context.Context, msg string) {
		fmt.Printf("Sending user msg to downstr: '%s'\n", msg)
		_, err := downstream.Write([]byte(msg))

		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
	})

	scanner := bufio.NewScanner(downstream)
	fmt.Printf("Starting to read server messages\n")
	for scanner.Scan() {
		in := scanner.Text()
		fmt.Printf("server: '%s'\n", in)
		serverMessaged.Emit(context.Background(), in)
	}
}
