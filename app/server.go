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

	header_map := make(map[string]string)
	for _, line := range request[2:] {
		header := strings.Split(line, ": ")

		if len(header) > 1 {
			header_map[header[0]] = header[1]	
		}
	}

	fmt.Println(header_map)

	// Default response
	response := []byte("HTTP/1.1 404 Not Found\r\n\r\n")

	// Handle different requests to HTTP server
	if path == "/" {
		response = []byte("HTTP/1.1 200 OK\r\n\r\n")
	} else {
		blocks := strings.Split(path, "/")[1:]
		fmt.Println(blocks)
		req_type := blocks[0]

		if req_type == "echo" {
			body := strings.Join(blocks[1:], "/")
			response = generateResponse(body)
		} else if req_type == "user-agent" {
			response = generateResponse(header_map["User-Agent"])
		}
	}

	conn.Write(response)

	// Close the connection when done
	conn.Close()
	l.Close()
}

func generateResponse(body string) []byte {
	response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(body), body)
	fmt.Println(response)
	return []byte(response)
}