package main

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/maniartech/signals"
)

type contextKey string

var userMessaged = signals.New[string]()
var serverMessaged = signals.New[string]()
var key contextKey = "id"
var tonysCoinAddr = "7YWHMfk9JZe0LM0g1ZauHuiSxhI"

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
			connection.Write([]byte(tonify(msg) + "\n"))
		}
	})

	downstream, err := net.Dial("tcp", "chat.protohackers.com:16963")
	fmt.Printf("Client started on port chat.protohackers.com:16963\n")

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Starting to read user messages\n")
	go handleServerConnection(downstream, ctx, cancel)

	buffer := make([]byte, 64*1024)
	for {
		n, err := connection.Read(buffer)
		in := strings.TrimSuffix(string(buffer[0:n]), "\n")
		fmt.Printf("user: '%s'\n", in)
		userMessaged.Emit(ctx, in)

		if err != nil {
			break
		}

		if ctx.Err() != nil {
			fmt.Printf("connection ctx has been cancelled\n")
			break
		}
	}

	cancel()

	// Tell downstream to stop listening on messages once user disconnected
	downstream.SetReadDeadline(time.Now())
}

func handleServerConnection(downstream net.Conn, ctx context.Context, cancel context.CancelFunc) {
	defer downstream.Close()
	defer cancel()

	userMessaged.AddListener(func(c context.Context, msg string) {
		if c.Value(key) == ctx.Value(key) {
			fmt.Printf("Sending user msg to downstr: '%s'\n", msg)
			downstream.Write([]byte(tonify(msg) + "\n"))
		}
	})

	buffer := make([]byte, 64*1024)
	for {
		n, err := downstream.Read(buffer)
		in := strings.TrimSuffix(string(buffer[0:n]), "\n")
		fmt.Printf("server: '%s'\n", in)
		serverMessaged.Emit(ctx, in)

		if err != nil {
			fmt.Printf("downstream read has been cancelled with time out\n")
			break
		}

		if ctx.Err() != nil {
			fmt.Printf("downstream ctx has been cancelled\n")
			break
		}
	}
}

func tonify(msg string) string {
	words := strings.Split(msg, " ")
	for _, w := range words {
		if len(w) > 0 && w[0] == '7' && len(w) >= 26 && len(w) <= 35 {
			// todo test for alphanum
			msg = strings.ReplaceAll(msg, w, tonysCoinAddr)
		}
	}

	return msg
}
