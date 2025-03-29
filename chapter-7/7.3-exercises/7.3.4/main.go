package main

import (
	"fmt"
	"sync"
	"time"

	"7.3.4/channel"
)

func main() {
	intChan := channel.NewChannel[int](2)
	wg := sync.WaitGroup{}
	wg.Add(1)

	go receiver(intChan, &wg)

	for i := 1; i <= 6; i++ {
		fmt.Println("Sending:", i)
		intChan.Send(i)
	}
	intChan.Send(-1)
	wg.Wait()
}

func receiver(messages *channel.Channel[int], wGroup *sync.WaitGroup) {
	msg := 0
	for msg != -1 {
		time.Sleep(1 * time.Second)
		msg = messages.Receive()
		fmt.Println("Received:", msg)
	}
	wGroup.Done()
}
