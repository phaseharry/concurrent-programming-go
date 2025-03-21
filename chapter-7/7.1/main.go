package main

import "fmt"

func main() {
	msgChannel := make(chan string)
	/*
	   example of an unbuffered channel. since we didn't state a capacity
	   for the channel, the channel can only hold one element. If it already
	   has 1 element, then any goroutine that attempts to send more / add more
	   messages to the channel will be blocked and must wait until a receiver
	   goroutine has consumed the message
	*/
	go receiver(msgChannel)

	fmt.Println("Sending HELLO to receiver go routine")
	msgChannel <- "HELLO"

	fmt.Println("Sending THERE to receiver go routine")
	msgChannel <- "THERE"

	fmt.Println("Sending STOP to receiver to have receiver goroutine stop iterating endlessly")
	msgChannel <- "STOP"
}

func receiver(messages chan string) {
	msg := ""
	for msg != "STOP" {
		msg = <-messages
		fmt.Println("Received:", msg)
	}
}
