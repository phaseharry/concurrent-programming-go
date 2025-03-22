package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	/*
	   creating a buffered channel of size 3.
	   goroutines can send items through it without being blocked
	   if the channel has not be filled up. Once full and a goroutine
	   attempts to send an element through then that goroutine will be blocked.
	*/
	msgChannel := make(chan int, 3)
	wGroup := sync.WaitGroup{}
	wGroup.Add(1)
	go receiver(msgChannel, &wGroup)

	/*
	   sending messages non-stop through channel. will fill up channel
	   by i == 3 and be blocked until the receiver consumes the messages
	   to reduce the size. this will unblock and let the main goroutine
	   send another element through the channel
	*/
	for i := 1; i <= 6; i++ {
		size := len(msgChannel)
		fmt.Printf("%s Sending: %d. Buffer Size: %d\n", time.Now().Format("15:04:05"), i, size)
		msgChannel <- i
	}
	// sending termination message to receiver goroutine
	msgChannel <- -1
	// waiting until the receiver goroutine is done running and calls wg.Done()
	wGroup.Wait()
}

func receiver(messages chan int, wGroup *sync.WaitGroup) {
	msg := 0
	/*
	   keeps iterating until we get a -1 indicating a termination fo this loop.
	   deliberately throttling the consume speed of receiver goroutine by
	   sleeping for a second before we consume a message. This will cause
	   the buffered messages channel to be filled up and block the main
	   goroutine until the receiver consumes messages to reduce the channel
	   size.
	*/
	for msg != -1 {
		time.Sleep(1 * time.Second)
		msg = <-messages
		fmt.Println("Received:", msg)
	}
	wGroup.Done()
}
