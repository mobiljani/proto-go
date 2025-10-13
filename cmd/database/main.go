package main

import (
	"fmt"
	"net"
	"strings"
	"sync"
)

type Store struct {
	mu sync.RWMutex
	db map[string]string
}

var store = Store{db: make(map[string]string)}

func (s *Store) add(key string, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.db[key] = value
}

func (s *Store) get(key string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.db[key]
}

func main() {
	// Resolve the string address to a UDP address
	udpAddr, err := net.ResolveUDPAddr("udp", ":8080")

	if err != nil {
		fmt.Println(err)
		return
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	fmt.Printf("Server started on %s\n", udpAddr.String())

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Read from UDP listener in endless loop
	buf := make([]byte, 32*1024)
	for {
		n, addr, err := conn.ReadFromUDP(buf[0:])
		if err != nil {
			fmt.Println(err)
			return
		}

		in := string(buf[0:n])

		fmt.Printf("Received '%s'\n", in)

		if strings.Contains(in, "version") {
			conn.WriteToUDP([]byte("version=Ken's Key-Value Store 1.0"), addr)
		} else if strings.Contains(in, "=") {
			key := strings.Split(in, "=")[0]
			value := strings.Replace(in, key+"=", "", 1)
			store.add(key, value)
		} else {
			m := fmt.Sprintf("%s=%s", in, store.get(in))
			conn.WriteToUDP([]byte(m), addr)
			fmt.Printf("Sent '%s'\n", string(m))
		}
	}

}
