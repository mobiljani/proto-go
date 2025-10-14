package main

import (
	"bufio"
	"context"
	"fmt"
	"net"

	"github.com/google/uuid"
	"github.com/maniartech/signals"
)

type contextKey string

var userMessaged = signals.New[string]()
var serverMessaged = signals.New[string]()
var key contextKey = "id"

func main() {

	list, err := net.Listen("tcp", ":8080")
	fmt.Printf("Server started on port 8080\n")

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	cb := context.Background()

	for {
		conn, err := list.Accept()

		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		ct := context.WithValue(cb, key, uuid.New().String())
		ctx, cancel := context.WithCancel(ct)
		go handleUserConnection(conn, ctx, cancel)

	}
}

func handleUserConnection(connection net.Conn, ctx context.Context, cancel context.CancelFunc) {
	defer connection.Close()
	defer cancel()

	serverMessaged.AddListener(func(c context.Context, msg string) {
		if c.Value(key) == ctx.Value(key) {
			fmt.Printf("Downstream msg to user: '%s'\n", msg)
			_, err := connection.Write([]byte(msg + "\n"))
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
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
	go handleServerConnection(downstream, ctx, cancel)

	for scanner.Scan() {
		in := scanner.Text()
		fmt.Printf("user: '%s'\n", in)
		userMessaged.Emit(ctx, in)

		if err := ctx.Err(); err != nil {
			break
		}
	}

}

func handleServerConnection(downstream net.Conn, ctx context.Context, cancel context.CancelFunc) {
	defer downstream.Close()
	defer cancel()

	userMessaged.AddListener(func(c context.Context, msg string) {
		if c.Value(key) == ctx.Value(key) {
			fmt.Printf("Sending user msg to downstr: '%s'\n", msg)
			_, err := downstream.Write([]byte(msg + "\n"))

			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
		}
	})

	scanner := bufio.NewScanner(downstream)
	fmt.Printf("Starting to read server messages\n")
	for scanner.Scan() {
		in := scanner.Text()
		fmt.Printf("server: '%s'\n", in)
		serverMessaged.Emit(ctx, in)

		if err := ctx.Err(); err != nil {
			break
		}
	}

}
