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
	
	request := strings.Split(string(buffer), "\r\n")
	start_line := strings.Split(request[0], " ")
	path := start_line[1]

	response := []byte("HTTP/1.1 404 Not Found\r\n\r\n")
	if path == "/" {
		response = []byte("HTTP/1.1 200 OK\r\n\r\n")
	} else {
		blocks := strings.Split(path, "/")
		if len(blocks) > 2 {
			req_type := blocks[1]
			req_str := strings.Join(blocks[2:], "/")

			if req_type == "echo" {
				data := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(req_str), req_str)
				response = []byte(data)
			}
		}
	}

	conn.Write(response)

	conn.Close()
	l.Close()
}
