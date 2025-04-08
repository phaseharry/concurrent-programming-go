package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	randomNumbers := generateNumbers()
	/*
	   when using time.After to create a timeout channel
	   outside of the select block, it will always send a timeout
	   time after the set time unlike if the channel was created within
	   a case.

	   If it was created within a case then that channel will only send
	   a message if no other cases have either sent or received a message
	   from any other channels within the cases
	*/
	timeout := time.After(5 * time.Second)
	for {
		select {
		case num := <-randomNumbers:
			fmt.Println("randomNumber:", num)
		case <-timeout:
			fmt.Println("5 seconds has passed")
			return
		}
	}
}

func generateNumbers() chan int {
	output := make(chan int)
	go func() {
		for {
			output <- rand.Intn(10)
			time.Sleep(200 * time.Millisecond)
		}
	}()
	return output
}
