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
	defer l.Close()
	
	// Run loop until no more connections to accept
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	_, err := conn.Read(buffer[:])
	if err != nil {
		fmt.Println("Error reading data: ", err.Error())
		os.Exit(1)
	}

	// Get initial request data
	request := strings.Split(string(buffer), "\r\n")
	start_line := strings.Split(request[0], " ")
	path := start_line[1]

	// Map headers for easy access
	header_map := make(map[string]string)
	for _, line := range request[2:] {
		header := strings.Split(line, ": ")

		if len(header) > 1 {
			header_map[header[0]] = header[1]	
		}
	}

	// Default response
	response := []byte("HTTP/1.1 404 Not Found\r\n\r\n")

	// Handle different requests to HTTP server
	if path == "/" {
		response = []byte("HTTP/1.1 200 OK\r\n\r\n")
	} else {
		blocks := strings.Split(path, "/")[1:]
		req_type := blocks[0]

		switch req_type {
		case "echo":
			body := strings.Join(blocks[1:], "/")
			response = generateResponse("text/plain", body)
		case "user-agent":
			response = generateResponse("text/plain", header_map["User-Agent"])
		case "files":
			filename := blocks[1]
			directory := os.Args[2]

			body, err := os.ReadFile(directory + filename)
			if err == nil {
				response = generateResponse("application/octet-stream", string(body))
			}
		}
	}
	
	// Write response to connection
	conn.Write(response)
}

func generateResponse(content_type string, body string) []byte {
	response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: %s\r\nContent-Length: %d\r\n\r\n%s", content_type, len(body), body)
	return []byte(response)
}