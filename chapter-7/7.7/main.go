package main

import (
	"fmt"
	"time"
)

func main() {
	msgChannel := make(chan int)
	go receiver(msgChannel)
	go sender(msgChannel)
	time.Sleep(5 * time.Second)
}

/*
adding the "<-" to the left side of the "chan" keyword will mark
this channel as a channel that consumes messages only. If you
try to send a message from this function then go will throw an
error at compile time
*/
func receiver(messages <-chan int) {
	for {
		msg := <-messages
		fmt.Println(time.Now().Format("15:04:05"), "Received:", msg)
	}
}

/*
adding the "<-" to the right side of the "chan" keyword will mark
this channel as a channel that only sends messages to that channel. If you
try to consume a message from this function then go will throw an
error at compile time.
*/
func sender(messages chan<- int) {
	for i := 1; ; i++ {
		fmt.Println(time.Now().Format("15:04:05"), "Sending:", i)
		messages <- i
		time.Sleep(1 * time.Second)
	}
}
