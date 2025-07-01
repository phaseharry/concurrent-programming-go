package main

import (
	"fmt"
	"net"
	"os"
	"regexp"
)

/*
Demo of the worker pool concurrency pattern.
Modified such that if all workers are busy processing and cannot
accept new tasks to process, it will send back a "Too Many Requests" message
back to client.
*/
func main() {
	incomingConnections := make(chan net.Conn)

	StartHttpWorkers(3, incomingConnections)

	server, _ := net.Listen("tcp", "localhost:8080")
	defer server.Close()
	for {
		conn, _ := server.Accept()
		select {
		case incomingConnections <- conn:
		default:
			fmt.Println("Server is busy")
			serverBusyResponse := "HTTP/1.1 429 Too Many Requests\r\n\r\n<html>Busy</html>\n"
			conn.Write([]byte(serverBusyResponse))
		}
	}
}

func StartHttpWorkers(n int, incomingConnections <-chan net.Conn) {
	for range n {
		go func() {
			for c := range incomingConnections {
				handleHttpRequest(c)
			}
		}()
	}
}

var r, _ = regexp.Compile("GET (.+) HTTP/1.1\r\n")

func handleHttpRequest(conn net.Conn) {
	buff := make([]byte, 1024)
	size, _ := conn.Read(buff)

	if r.Match(buff[:size]) {
		file, err := os.ReadFile(
			fmt.Sprintf("../resources/%s", r.FindSubmatch(buff[:size])[1]),
		)

		if err == nil {
			response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Length: %d\r\n\r\n", len(file))
			conn.Write([]byte(response))
			conn.Write(file)
		} else {
			notFoundResponse := "HTTP/1.1 404 Not Found\r\n\r\n<html>Not Found</html>"
			conn.Write([]byte(notFoundResponse))
		}
	} else {
		internalServerErrResponse := "HTTP/1.1 500 Internal Server Error\r\n\r\n"
		conn.Write([]byte(internalServerErrResponse))
	}
	conn.Close()
}
