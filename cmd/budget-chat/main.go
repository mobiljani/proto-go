package main

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/maniartech/signals"
)

type Message struct {
	user string
	msg  string
}

var userJoined = signals.New[string]()
var userLeft = signals.New[string]()
var messageSent = signals.New[Message]()
var names []string

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
	var name string

	defer connection.Close()

	userJoined.AddListener(func(ctx context.Context, user string) {
		m := fmt.Sprintf("* %s has entered the room\n", user)
		connection.Write([]byte(m))
	})

	userLeft.AddListener(func(ctx context.Context, user string) {
		m := fmt.Sprintf("* %s has left the room\n", user)
		connection.Write([]byte(m))
	})

	messageSent.AddListener(func(ctx context.Context, message Message) {
		if name != "" && name != message.user {
			m := fmt.Sprintf("[%s] %s\n", message.user, message.msg)
			connection.Write([]byte(m))
		}
	})

	fmt.Printf("Client connected: %s\n", connection.RemoteAddr().String())
	connection.Write([]byte("Welcome to budgetchat! What shall I call you?\n"))
	scanner := bufio.NewScanner(connection)

	for scanner.Scan() {
		rb := scanner.Bytes()
		in := string(rb)
		fmt.Printf("Received: %s \n", rb)

		if name == "" {
			if !validateName(in) {
				connection.Write([]byte("Name must be between 1 and 16 characters\n"))
				break
			}

			name = in

			if len(names) > 0 {
				m := fmt.Sprintf("* The room contains: %s\n", strings.Join(names, ", "))
				connection.Write([]byte(m))
			}

			names = append(names, in)

			ctx := context.Background()
			userJoined.Emit(ctx, name)

			continue
		}

		ctx := context.Background()
		messageSent.Emit(ctx, Message{user: name, msg: in})
	}

	if name != "" {
		ctx := context.Background()
		userLeft.Emit(ctx, name)
	}

}

func validateName(name string) bool {
	return len(name) >= 1 && len(name) <= 16
}
