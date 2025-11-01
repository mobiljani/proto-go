package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {

	url := "vcs.protohackers.com"
	port := 30307

	add, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", url, port))
	conn, _ := net.DialTCP("tcp", nil, add)

	r := bufio.NewReader(conn)

	// PUT
	read(r)
	send("PUT /test.txt 8\n", conn)
	send("aaaaaaa\n", conn)
	read(r)

	//GET
	send("GET /test.txt\n", conn)
	read(r)
	read(r)
	read(r)
	read(r)
	read(r)
	read(r)
}

func send(cmd string, conn *net.TCPConn) {
	content := []byte(cmd)
	_, _ = conn.Write(content)
	fmt.Printf("Sending: %s", cmd)
}

func read(r *bufio.Reader) {
	reply, _ := r.ReadBytes(byte('\n'))
	fmt.Printf("Received: %s", string(reply))
}
