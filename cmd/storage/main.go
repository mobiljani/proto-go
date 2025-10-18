package main

import (
	"bufio"
	"fmt"
	"net"
	"regexp"
	"strings"
	"time"
)

func main() {
	port := 8080
	fmt.Printf("Starting server on port %d\n", port)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		fmt.Printf("Failed to create a server: %v\n", err)
		return
	}

	for {
		connection, err := listener.Accept()

		if err != nil {
			fmt.Printf("Connection error %v\n", err)
			continue
		}

		go handleConnection(connection)
	}

}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Printf("Connection from %s\n", conn.RemoteAddr().String())
	say(conn, "READY\n")

	r := bufio.NewReader(conn)

	for {
		line, err := r.ReadString('\n')
		if err != nil {
			break
		}

		switch cmd := getCmd(line); cmd {
		case "HELP":
			say(conn, "OK usage: HELP|GET|PUT|LIST\n")
		case "LIST":
			list(conn, line)
		case "GET":
			get(conn, line)
		case "PUT":
			put(conn, line)
		default:
			fmt.Printf("ERR illegal method: '%s'", cmd)
			return
		}
		say(conn, "READY\n")
	}
}

var filename = regexp.MustCompile(`^/[a-zA-Z0-9] [ ._-]{2,256}$`)

func list(conn net.Conn, line string) {
	params := getParams(line)

	if len(params) != 1 || len(params) != 2 {
		say(conn, "ERR usage: GET file [revision]\n")
		return
	}

	if !filename.Match([]byte(params[0])) {
		say(conn, "ERR illegal file name\n")
		return
	}

	// TODO continue here

}

func get(conn net.Conn, line string) {

}

func put(conn net.Conn, line string) {

}

func getCmd(line string) string {
	cmd := line
	if strings.Contains(line, " ") {
		cmd = line[0:strings.Index(line, " ")]
	}

	if strings.HasSuffix(line, "\n") {
		cmd = line[0 : len(cmd)-1]
	}

	return strings.ToUpper(cmd)
}

func getParams(line string) []string {
	params := line

	if strings.HasSuffix(line, "\n") {
		params = line[0 : len(params)-1]
	}

	if strings.HasPrefix(params, " ") || strings.HasSuffix(params, " ") {
		params = strings.Trim(params, " ")
	}

	return strings.Split(params, " ")
}

func say(conn net.Conn, msg string) {
	for _, c := range msg {
		conn.Write([]byte(string(c)))
		time.Sleep(25 * time.Millisecond)
	}
}
