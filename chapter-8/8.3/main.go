package main

import (
	"fmt"
	"time"
)

func main() {
	messages := sendMsgAfter(3 * time.Second)
	for {
		/*
		   normally, using the select case blocks when there are no messages from
		   your list of channels to consume from. To make it non-blocking, a default
		   case can be added and it will process that go block instead when there are no messages
		   available for other cases
		*/
		select {
		case msg := <-messages:
			fmt.Println("Message received:", msg)
			return
		default:
			fmt.Println("No messages waiting")
			time.Sleep(1 * time.Second)
		}
	}
}

func sendMsgAfter(seconds time.Duration) <-chan string {
	messages := make(chan string)
	go func() {
		time.Sleep(seconds)
		messages <- "Hello"
	}()
	return messages
}
