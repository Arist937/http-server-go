package main

import (
	"fmt"
	"net"
	"os"
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
	conn.Read(buffer[:])

	data := []byte("HTTP/1.1 200 OK\r\n\r\n")
	conn.Write(data)
}
