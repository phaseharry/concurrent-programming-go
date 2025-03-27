package main

import "fmt"

func main() {
	msgChannel := make(chan string)
	go receiver(msgChannel)

	fmt.Println("Sending HELLO to receiver go routine")
	msgChannel <- "HELLO"

	fmt.Println("Sending THERE to receiver go routine")
	msgChannel <- "THERE"

	fmt.Println("Sending STOP to receiver to have receiver goroutine stop iterating endlessly")
	msgChannel <- "STOP"

	/*
	   attempting to consume a message in main goroutine. will block until
	   there's a message to be able to be consumed. this will prevent the main goroutine
	   from terminating before we print the "STOP" message
	*/
	<-msgChannel
	close(msgChannel)
}

func receiver(messages chan string) {
	for msg := range messages {
		fmt.Println("Received:", msg)
		// if we get "STOP" break out of the loop for consuming messages,
		if msg == "STOP" {
			break
		}
	}
	/*
	   after we break out of message consuming loop when we get the "STOP" message,
	   we send an arbitrary message through to the channel so that the main
	   goroutine has something to consume and be unblocked.
	*/
	messages <- ""
}
