package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {	
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	buffer := make([]byte, 1024)
	_, err = conn.Read(buffer[:])
	if err != nil {
		fmt.Println("Error reading data: ", err.Error())
		os.Exit(1)
	}
	
	stringified := string(buffer)
	fields := strings.Fields(stringified)
	path := fields[1]

	if path == "/" {
		data := []byte("HTTP/1.1 200 OK\r\n\r\n")
		conn.Write(data)
	} else {
		data := []byte("HTTP/1.1 404 Not Found\r\n\r\n")
		conn.Write(data)
	}
}
