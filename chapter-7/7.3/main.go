package main

import (
	"fmt"
	"time"
)

func main() {
	msgChannel := make(chan string)
	go deadReceiver(msgChannel)

	/*
	   because we're using an unbuffered msgChannel,
	   when we send a message to a channel, there must be
	   a goroutine that consumes that message before
	   the goroutine that sent it can process any further.
	   in this case, the main goroutine is blocked
	   after sending "No one will hear this" because
	   the deadReceiver goroutine does not consume the message,
	   leading to a deadlock.
	*/
	msgChannel <- "No one will hear this"

}

func deadReceiver(messages chan string) {
	time.Sleep(5 * time.Second)
	fmt.Println("Receiver slept for 5 seconds")
}
