package main

import (
	"fmt"
	"time"
)

func main() {
	messagesFromA := writeEvery("Tick", 1*time.Second)
	messagesFromB := writeEvery("Tock", 3*time.Second)

	for {
		/*
		   creating an infinite loop here that will continuous block and listen
		   for messages from the 2 open channels. if there's a message for channelA
		   then the first case will be processed and you can have special logic just
		   used for the messages from channelA, same for the channelB.

		   if there's a message that's being received from multiple channels at the same time,
		   a message is chosen at random in which the order of processing the cases is done.
		   code should not rely on a specific order of case processing to work correctly.
		*/
		select {
		case msg1 := <-messagesFromA:
			fmt.Println(msg1)
		case msg2 := <-messagesFromB:
			fmt.Println(msg2)
		}
	}
}

func writeEvery(msg string, seconds time.Duration) <-chan string {
	messages := make(chan string)

	/*
	   creates a new messages channel and returns it. also triggers a goroutine
	   that will continously sleep for x seconds and then attempt to consume a message.
	*/
	go func() {
		for {
			time.Sleep(seconds)
			messages <- msg
		}
	}()

	return messages
}
