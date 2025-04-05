package main

import "fmt"

func main() {
	var ch chan string = nil
	/*
	   when a channel is set to nil (it's default value if it wasn't created with make())
	   any attempts to send a messages through or receive a message from that channel
	   will be blocked indefinitely.
	*/
	ch <- "message"
	fmt.Println("This will never be printed")
}
