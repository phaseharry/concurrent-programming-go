package main

import (
	"fmt"
	"net"
	"os"
	"regexp"
)

/*
Demo of the worker pool concurrency pattern.
In this pattern, n amount of parallel threads / goroutines are created initially
and each of the workers will process tasks through a common queue. This ensure that
tasks can be processed in parallel while limiting the amount of memory used since there
can be a max of n threads / goroutines. Threads / goroutines can be increased if the load of
tasks reach a threshold, etc. If there's no work then the spun up workers will idle until there is.
This pattern is useful in an environment where spinning up new threads of execution dynamically is
expensive since there's only an initial start up. Not as useful with regard to performance in go
since goroutines are not tied to a thread 1:1 (m:n threads. 1 thread of CPU execution can be responsible for multiple goroutines).
Ex: Having a Http Server consume messages and having 3 workers processing the messages.
*/
func main() {
	incomingConnections := make(chan net.Conn)

	// spin up 3 web workers that will consume connections and process HTTP requests
	StartHttpWorkers(3, incomingConnections)

	server, _ := net.Listen("tcp", "localhost:8080")
	defer server.Close()
	for {
		// blocks until there's a connection and then enqueue so a worker can process request	0
		conn, _ := server.Accept()
		incomingConnections <- conn
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
	conn.Close()
}
