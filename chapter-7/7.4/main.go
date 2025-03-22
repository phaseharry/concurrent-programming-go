package main

import (
	"fmt"
	"time"
)

func main() {
	msgChannel := make(chan string)
	go deadSender(msgChannel)
	fmt.Println("Reading message from channel...")

	/*
		msgChannel is not a buffered channel so consume / send calls will block the
		goroutines that calls them and waits until the message has been consumed / received.
		in this case, the main goroutine will be blocked because it's trying to consume
		from an empty channel because the deadSender goroutine is not sending any messages
		through the channel to be consumed.
	*/
	msg := <-msgChannel
	fmt.Println("Received:", msg)
}

func deadSender(messages chan string) {
	time.Sleep(5 * time.Second)
	fmt.Println("Sender slept for 5 seconds")
}
