package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func handlePost(reader *bufio.Reader, filename string, dir string) string {
	fmt.Println("handlePost: ", filename)
	fmt.Println("Directory: ", dir)

	filepath := filepath.Join(dir, filename)

	// Read the headers
	headers := make(map[string]string)
	for {
		line, err := reader.ReadString('\n')
		if err != nil || line == "\r\n" {
			break
		}

		parts := strings.SplitN(line, ": ", 2)
		if len(parts) == 2 {
			headers[parts[0]] = strings.TrimSpace(parts[1])
		}
	}

	// Get the Content-Length header
	contentLengthStr, ok := headers["Content-Length"]
	if !ok {
		return "HTTP/1.1 400 Bad Request\r\n\r\n"
	}

	fmt.Println("contentLengthStr: ", contentLengthStr)

	// Convert the Content-Length to an integer
	contentLength, err := strconv.Atoi(contentLengthStr)
	if err != nil {
		return "HTTP/1.1 400 Bad Request\r\n\r\n"
	}

	fmt.Println("contentLength: ", contentLength)

	// Create a limited reader and read the body
	body, err := ioutil.ReadAll(io.LimitReader(reader, int64(contentLength)))
	if err != nil {
		return "HTTP/1.1 500 Internal Server Error\r\n\r\n"
	}

	fmt.Println("body: ", string(body))

	// Write the body to the file
	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return "HTTP/1.1 500 Internal Server Error\r\n\r\n"
	}
	defer file.Close()
	file.Write(body)

	return "HTTP/1.1 201 Created\r\n\r\n"
}

func getAndHandleFiles(conn net.Conn, filename string, dir string) string {
	filepath := filepath.Join(dir, filename)
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return "HTTP/1.1 404 Not Found\r\n\r\n"
	}
	content, err := os.ReadFile(filepath)
	if err != nil {
		return "HTTP/1.1 404 Not Found\r\n\r\n"
	}

	return fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s\r\n", len(content), string(content))
}

func handleConnection(conn net.Conn, dir string) {
	defer conn.Close()

	headerType := ""
	var responseBody string
	var response string
	idx := 0
	reader := bufio.NewReader(conn)

loop:
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
			header := strings.Split(line, " ")
			method := header[0]
			path := header[1]
			command := strings.Split(path, "/")[1]

			if command == "echo" {
				headerType = "echo"
				responseBody = strings.TrimPrefix(path, "/echo/")
				fmt.Printf("echo responseBody: %s\n", responseBody)
				break
			} else if command == "user-agent" {
				headerType = "user-agent"
			} else if command == "files" {
				headerType = "files"
				switch method {
				case "GET":
					filename := strings.TrimPrefix(path, "/files/")
					response = getAndHandleFiles(conn, filename, dir)
				case "POST":
					filename := strings.TrimPrefix(path, "/files/")
					response = handlePost(reader, filename, dir)
					break loop
				}
			} else {
				headerType = path
			}
		} else {
			if headerType == "user-agent" && strings.HasPrefix(line, "User-Agent: ") {
				responseBody = strings.TrimSpace(strings.TrimPrefix(line, "User-Agent: "))
				fmt.Printf("User-Agent: responseBody: %s\n", responseBody)
				fmt.Printf("User-Agent len: %d\n", len(responseBody))
			}
		}
	}

	switch headerType {
	case "/":
		fmt.Println("No header")
		response = fmt.Sprintf("HTTP/1.1 200 OK\r\n\r\n")
	case "echo":
		fmt.Println("echo")
		response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(responseBody), responseBody)
	case "user-agent":
		response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(responseBody), responseBody)
	case "files":
		break
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
	dirFlag := flag.String("directory", ".", "the directory of static file to host")
	flag.Parse()

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConnection(conn, *dirFlag)
	}
}
