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
	for {
		msg, channelActive := <-messages
		fmt.Println(time.Now().Format("15:04:05"), "Received:", msg, channelActive)
		/*
		   when consuming messages from channels, there's a second value that's
		   a flag value, indicating whether the channel that we're consuming messages
		   from is closed or not. if it's closed / not active, that value will be false,
		   otherwise it is true
		*/
		if !channelActive {
			return
		}
	}
}
