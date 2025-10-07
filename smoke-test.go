package main

import (
	"fmt"
	"net"
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

func handleConnection(connection net.Conn) {
	defer connection.Close()

	// _, err := io.Copy(connection, connection)
	// if err != nil {
	// 	fmt.Printf("Error during io.Copy: %v\n", err)
	// }

	fmt.Printf("Connection from %s\n", connection.RemoteAddr().String())

	buffer := make([]byte, 32*1024)

	for {
		bytesRead, readError := connection.Read(buffer)
		fmt.Printf("Read %d bytes\n", bytesRead)
		if bytesRead > 0 {

			bytesWritten, writeError := connection.Write(buffer[0:bytesRead])
			fmt.Printf("Wrote %d bytes\n", bytesWritten)

			if writeError != nil {
				fmt.Printf("Error during write: %v\n", writeError)
				break
			}

			if bytesWritten < 0 || bytesWritten > bytesRead {
				fmt.Printf("Wrote invalid number of bytes: %d\n", bytesWritten)
				break
			}
		}
		if readError != nil {
			fmt.Printf("Error during read: %v\n", readError)
			break
		}
	}

}
