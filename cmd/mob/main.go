package main

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"strings"

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
			connection.Write([]byte(msg + "\n"))
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

		if ctx.Err() != nil {
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
			downstream.Write([]byte(tonify(msg) + "\n"))
		}
	})

	scanner := bufio.NewScanner(downstream)
	fmt.Printf("Starting to read server messages\n")
	for scanner.Scan() {
		in := scanner.Text()
		fmt.Printf("server: '%s'\n", in)
		serverMessaged.Emit(ctx, in)

		if ctx.Err() != nil {
			break
		}
	}

}

func tonify(msg string) string {
	words := strings.Split(msg, " ")

	for _, w := range words {
		if w[0] == '7' && len(w) >= 26 && len(w) <= 35 {
			// todo test for alphanum
			msg = strings.ReplaceAll(msg, w, tonysCoinAddr)
		}
	}

	return msg

}
