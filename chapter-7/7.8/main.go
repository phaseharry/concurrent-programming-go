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
	/*
				closing the msgChannel after 6 seconds.
		    - send i = 1, sleep for 1 second
		    - send i = 2, sleep for 1 second
		    - send i = 3, sleep for 1 second
		    - close the channel and sleep for 3 seconds
				after a channel is closed, goroutines should not be sending
				messages to it as it will raise errors.
				goroutines that tries to consume messages from a closed channel
				will receive the default value for the type of the channel elements.
	*/
	close(msgChannel)
	time.Sleep(3 * time.Second)
}

func receiver(messages <-chan int) {
	for {
		msg := <-messages
		fmt.Println(time.Now().Format("15:04:05"), "Received:", msg)
		/*
		   consume a message then wait 1 second before consuming another
		   message. After consuming the initial i values (1-3), the
		   channel is closed in the main goroutine.
		   the receiver goroutine will only able to get the default 0 value of type int
		*/
		time.Sleep(1 * time.Second)
	}
}
