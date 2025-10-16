package main

import (
	"bufio"
	"fmt"
	"net"
	"regexp"
	"strings"
)

var tonysCoinAddr = "7YWHMfk9JZe0LM0g1ZauHuiSxhI"

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
		go handleUpstreamConnection(conn)
	}
}

func handleUpstreamConnection(upstream net.Conn) {
	defer upstream.Close()

	downstream, err := net.Dial("tcp", "chat.protohackers.com:16963")
	fmt.Printf("Client started on port chat.protohackers.com:16963\n")

	defer downstream.Close()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Starting to read user messages\n")
	go handleDownstreamConnection(upstream, downstream)

	r := bufio.NewReader(upstream)
	for in, err := r.ReadString('\n'); ; {
		if err != nil {
			return
		}
		fmt.Printf("user: '%s'\n", in)
		downstream.Write([]byte(tonify(in) + "\n"))
	}
}

func handleDownstreamConnection(upstream net.Conn, downstream net.Conn) {
	defer downstream.Close()
	defer upstream.Close()

	r := bufio.NewReader(downstream)
	for in, err := r.ReadString('\n'); ; {
		if err != nil {
			return
		}
		fmt.Printf("server: '%s'\n", in)
		upstream.Write([]byte(tonify(in) + "\n"))
	}
}

func tonify(msg string) string {
	words := strings.Split(msg, " ")
	for _, w := range words {
		if len(w) > 0 && w[0] == '7' && len(w) >= 26 && len(w) <= 35 {
			is_alphanumeric := regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(w)
			if is_alphanumeric {
				msg = strings.ReplaceAll(msg, w, tonysCoinAddr)
			}
		}
	}

	return msg
}
