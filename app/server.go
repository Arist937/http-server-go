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

	buffer := make([]byte, 2048)
	bytes, err := conn.Read(buffer[:])
	if err != nil {
		fmt.Println("Error reading data: ", err.Error())
		os.Exit(1)
	}

	// Get initial request data
	request := strings.Split(string(buffer[:bytes]), "\r\n")
	start_line := strings.Split(request[0], " ")
	method := start_line[0]
	path := start_line[1]

	// Map headers and get body if it exists
	header_map := make(map[string]string)

	index := 2
	for request[index] != "" {
		header := strings.Split(request[index], ": ")
		header_map[header[0]] = header[1]	

		index++
	}

	body := ""
	if index + 1 < len(request) {
		body = strings.TrimSpace(request[index + 1])
	}

	// Default response
	response := []byte("HTTP/1.1 404 Not Found\r\n\r\n")

	switch method {
	case "GET":
		response = handleGET(path, header_map)
	case "POST":
		response = handlePOST(path, body)
	}
	
	// Write response to connection
	conn.Write(response)
}

func handlePOST(path string, body string) []byte {
	blocks := strings.Split(path, "/")[1:]
	req_type := blocks[0]

	switch req_type {
	case "files":
		filename := blocks[1]
		directory := os.Args[2]
		os.WriteFile(directory + filename, []byte(body), os.FileMode(0666))

		return []byte("HTTP/1.1 201 OK\r\n\r\n")
	}

	return []byte("HTTP/1.1 404 Not Found\r\n\r\n")
}

func handleGET(path string, header_map map[string]string) []byte {
	if path == "/" {
		return []byte("HTTP/1.1 200 OK\r\n\r\n")
	} else {
		blocks := strings.Split(path, "/")[1:]
		req_type := blocks[0]

		switch req_type {
		case "echo":
			body := strings.Join(blocks[1:], "/")
			return generateResponse("text/plain", body)
		case "user-agent":
			return generateResponse("text/plain", header_map["User-Agent"])
		case "files":
			filename := blocks[1]
			directory := os.Args[2]
			body, err := os.ReadFile(directory + filename)

			if err == nil {
				return generateResponse("application/octet-stream", string(body))
			}
		}
	}

	return []byte("HTTP/1.1 404 Not Found\r\n\r\n")
}

func generateResponse(content_type string, body string) []byte {
	response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: %s\r\nContent-Length: %d\r\n\r\n%s", content_type, len(body), body)
	return []byte(response)
}