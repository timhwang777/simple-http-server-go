# Simple HTTP Server in Golang

![Static Badge](https://img.shields.io/badge/Go-Solutions-blue?logo=Go
)

## Table of Contents
1. [About the Project](#about-the-project)
2. [Getting Started](#getting-started)
    - [Prerequisites](#prerequisites)
3. [Author](#author)

## About the Project

This project is a simple HTTP server written in Go. It listens for incoming connections and handles them concurrently. The server supports three commands: `echo`, `user-agent`, and `files`. 

- The `echo` command echoes back the path after the `/echo/` in the response body.
- The `user-agent` command returns the User-Agent from the request headers.
- The `files` command supports both GET and POST methods. The GET method reads a file from the server's directory and returns its content. The POST method receives a file from the client and saves it to the server's directory.

The server is designed to be simple and easy to understand, making it a great starting point for anyone interested in learning about network programming in Go.

## Getting Started

To get a local copy up and running, follow these simple steps:

### Prerequisites

- Go: You need to have Go installed on your machine to run this server. You can download it from the official website: [https://golang.org/dl/](https://golang.org/dl/)

To run the server, navigate to the directory containing the `server.go` file in your terminal and run the command `go run server.go`.

Alternatively, you can execute the `your_server.sh` shell script.

## Author
Timothy Hwang