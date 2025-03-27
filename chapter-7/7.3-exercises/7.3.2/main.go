package main

import (
	"fmt"
	"time"
)

func main() {
	stringChannel := make(chan string)
	intSliceChannel := make(chan []int)
	go stringReceiver(stringChannel)
	go intSliceReceiver(intSliceChannel)
	for i := 1; i <= 3; i++ {
		fmt.Println(time.Now().Format("15:04:05"), "Sending:", i)
		stringChannel <- fmt.Sprintf("hello %d", i)
		intSliceChannel <- []int{i}
		time.Sleep(1 * time.Second)
	}
	close(stringChannel)
	close(intSliceChannel)
	time.Sleep(3 * time.Second)
}

func stringReceiver(messages <-chan string) {
	for {
		msg := <-messages
		fmt.Println(time.Now().Format("15:04:05"), "Received:", msg)
		time.Sleep(1 * time.Second)
	}
}

func intSliceReceiver(messages <-chan []int) {
	for {
		msg := <-messages
		fmt.Println(time.Now().Format("15:04:05"), "Received:", msg)
		time.Sleep(1 * time.Second)
	}
}
