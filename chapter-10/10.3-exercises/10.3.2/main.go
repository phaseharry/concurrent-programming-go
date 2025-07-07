package main

import (
	"fmt"
	"net"
	"os"
	"regexp"
)

func main() {
	/*
	   Set work queue to have a buffer of 10 so if all worker goroutines
	   are busy processing, we'll still listen to up to 10 incoming requests
	   and enqueue them before we send 503 Service Unavailable messages instead
	*/
	incomingConnections := make(chan net.Conn, 10)

	// spin up 3 web workers that will consume connections and process HTTP requests
	StartHttpWorkers(3, incomingConnections)

	server, _ := net.Listen("tcp", "localhost:8080")
	defer server.Close()
	for {
		// blocks until there's a connection and then enqueue so a worker can process request	0
		conn, _ := server.Accept()

		select {
		case incomingConnections <- conn:
		default:
			fmt.Println("Server is busy")
			conn.Write([]byte("HTTP/1.1 503 Service Unavailable\r\n\r\n" + "<html>Busy</html>\n"))
			conn.Close()
		}
	}
}

/*
initializes n workers that will consume connections from a common channel
that acts as a queue to enqueue messages for the workers to process.
*/
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

/*
reads in a http connection and response with the requested content.
if request is valid, attempt to respond with file if it's available, otherwise return not found.
if request is not valid, return 500 internal server error
*/
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
}
