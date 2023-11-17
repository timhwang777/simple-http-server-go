package main

import (
	"bufio"
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

	defer conn.Close()

	readBuf := make([]byte, 1024)
	_, err = conn.Read(readBuf)
	if err != nil {
		fmt.Println("Error reading: ", err.Error())
		os.Exit(1)
	}

	headerType := ""
	var responseBody string
	reader := bufio.NewReader(conn)
	for {
		fmt.Println("Reading...")
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		if line == "\r\n" {
			break
		}
		if strings.HasPrefix(line, "User-Agent: ") {
			headerType = "User-Agent: "
			responseBody = strings.TrimPrefix(line, "User-Agent: ")
		}
		if strings.HasPrefix(line, "/echo/") {
			headerType = "echo"
			responseBody = strings.TrimPrefix(line, "/echo/")
		}
		fmt.Println("Bottom Reading...")
	}
	fmt.Println("End reading...")

	var response string
	switch headerType {
	case "":
		response = fmt.Sprintf("HTTP/1.1 200 OK\r\n\r\n")
	case "echo":
		response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(responseBody), responseBody)
	case "User-Agent: ":
		response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(responseBody), responseBody)
	default:
		response = fmt.Sprintf("HTTP/1.1 404 Not Found\r\n\r\n")
	}

	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing: ", err.Error())
		os.Exit(1)
	}
}
