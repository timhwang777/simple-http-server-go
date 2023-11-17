package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	headerType := ""
	var responseBody string
	idx := 0
	reader := bufio.NewReader(conn)
	for {
		idx++
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		if line == "\r\n" {
			break
		}
		// first line for header
		if idx == 1 {
			header := strings.Split(line, " ")[1]
			fmt.Println("header: ", header)
			command := strings.Split(header, "/")[1]
			fmt.Println("command: ", command)

			if command == "echo" {
				headerType = "echo"
				responseBody = strings.TrimPrefix(header, "/echo/")
				fmt.Printf("echo responseBody: %s\n", responseBody)
				break
			} else if command == "user-agent" {
				headerType = "user-agent"
			} else {
				headerType = header
			}
		} else {
			if headerType == "user-agent" && strings.HasPrefix(line, "User-Agent: ") {
				responseBody = strings.TrimSpace(strings.TrimPrefix(line, "User-Agent: "))
				fmt.Printf("User-Agent: responseBody: %s\n", responseBody)
				fmt.Printf("User-Agent len: %d\n", len(responseBody))
			}
		}
	}

	var response string
	switch headerType {
	case "/":
		fmt.Println("No header")
		response = fmt.Sprintf("HTTP/1.1 200 OK\r\n\r\n")
	case "echo":
		fmt.Println("echo")
		response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(responseBody), responseBody)
	case "user-agent":
		response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(responseBody), responseBody)
	default:
		response = fmt.Sprintf("HTTP/1.1 404 Not Found\r\n\r\n")
	}

	fmt.Println("response: ", response)
	_, err := conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing: ", err.Error())
		os.Exit(1)
	}
}

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
	go handleConnection(conn)
}
