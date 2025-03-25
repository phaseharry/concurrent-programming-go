package main

import (
	"fmt"
	"time"
)

func main() {
	msgChannel := make(chan int)
	go receiver(msgChannel)
	for i := 1; i <= 3; i++ {
		fmt.Println(time.Now().Format("15:04:05"), "Sending:", i)
		msgChannel <- i
		time.Sleep(1 * time.Second)
	}
	close(msgChannel)
	time.Sleep(3 * time.Second)
}

func receiver(messages <-chan int) {
	/*
	   a cleaner approach to consume messages until the channel has been closed.
	   when using range, we will continue iterating and consuming messages until
	   the channel has been closed. if the channel is empty then the receiver
	   goroutine will block until there is a message to consume
	*/
	for msg := range messages {
		fmt.Println(time.Now().Format("15:04:05"), "Received:", msg)
		time.Sleep(1 * time.Second)
	}
	fmt.Println("Receiver finished.")
}
