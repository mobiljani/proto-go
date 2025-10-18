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

		downstream, err := net.Dial("tcp", "chat.protohackers.com:16963")
		fmt.Printf("Client started on port chat.protohackers.com:16963\n")
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		go handleUpstreamConnection(conn, downstream)
	}
}

func handleUpstreamConnection(upstream net.Conn, downstream net.Conn) {
	defer upstream.Close()
	defer downstream.Close()

	fmt.Printf("Starting to read user messages\n")
	go handleDownstreamConnection(upstream, downstream)

	relay(upstream, downstream)
}

func handleDownstreamConnection(upstream net.Conn, downstream net.Conn) {
	defer downstream.Close()
	defer upstream.Close()

	relay(downstream, upstream)
}

func relay(from net.Conn, to net.Conn) {
	r := bufio.NewReader(from)
	for {
		in, err := read(r)
		if err != nil {
			return
		}
		fmt.Printf("from: '%s'\n", in)

		_, err = to.Write([]byte(tonify(in)))
		if err != nil {
			fmt.Print("write error\n", err)
			return
		}
	}
}

func read(r *bufio.Reader) (string, error) {
	line, err := r.ReadString('\n')

	if err == nil {
		if line[len(line)-1] == '\n' {
			line = line[:len(line)-1]
		}
	}

	return line, err
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

	return msg + "\n"
}
