package main

import (
	"bufio"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type FileRecord struct {
	path    string
	version uint16
	lenght  uint64
	data    []byte
}

var files = []FileRecord{}

func main() {
	port := 9876
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
		fmt.Println(line)
		switch cmd := getCmd(line); cmd {
		case "HELP":
			say(conn, "OK usage: HELP|GET|PUT|LIST\n")
		case "LIST":
			list(conn, line)
		case "GET":
			get(conn, line)
		case "PUT":
			put(conn, line, r)
		default:
			fmt.Printf("ERR illegal method: '%s'", cmd)
			return
		}
		say(conn, "READY\n")
	}
}

var filename = regexp.MustCompile(`^/[a-zA-Z0-9._-]{2,256}$`)

func list(conn net.Conn, line string) {
	params := getParams(line)

	if len(params) != 1 {
		say(conn, "ERR usage: GET file [revision]\n")
		return
	}

	if !filename.Match([]byte(params[1])) {
		say(conn, "ERR illegal file name\n")
		return
	}

	file := params[1]
	c := 0

	for _, f := range files {
		if f.path == file {
			say(conn, fmt.Sprintf("Found %s %s %d", f.path, f.version, f.lenght))
			c += 1
		}
	}

	say(conn, fmt.Sprintf("OK %d", c))
}

func get(conn net.Conn, line string) {
	params := getParams(line)

	if len(params) <= 1 && len(params) > 3 {
		say(conn, "ERR usage: GET file [revision]\n")
		return
	}

	if !filename.Match([]byte(params[1])) {
		say(conn, "ERR illegal file name\n")
		return
	}

	file := params[1]

	for _, f := range files {
		if f.path == file {
			fmt.Printf("Found %s %d %d\n", f.path, f.version, f.lenght)
			say(conn, "READY\n")
			say(conn, fmt.Sprintf("OK %d\n", f.lenght))
			conn.Write(f.data)

			return

		}
	}

	say(conn, "ERR no such file\n")

}

func put(conn net.Conn, line string, r *bufio.Reader) {
	params := getParams(line)

	if len(params) != 3 {
		say(conn, "ERR usage: PUT file length newline data\n")
		return
	}

	if !filename.Match([]byte(params[1])) {
		say(conn, "ERR illegal file name\n")
		return
	}

	file := params[1]
	len := params[2]
	leni, _ := strconv.ParseUint(len, 10, 64)

	data, err := r.ReadBytes(byte('\n'))
	if err != nil {
		fmt.Print(err)
	}
	fmt.Printf("data lenght %d\n", cap(data))

	//todo look up version

	f := FileRecord{path: file, lenght: leni, version: 1, data: data}

	files = append(files, f)

	say(conn, "OK r1\n")

}

func getCmd(line string) string {
	cmd := line
	if strings.Contains(line, " ") {
		cmd = line[0:strings.Index(line, " ")]
	}

	if strings.HasSuffix(cmd, "\n") {
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
		time.Sleep(5 * time.Millisecond)
	}
	fmt.Println(msg)
}
