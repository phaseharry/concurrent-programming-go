package main

import (
	"fmt"
	"math/rand"
	"time"
)

/*
only outputing the latest temp when the outputTemp goroutine
is ready during each 2 second interval.
*/
func main() {
	temps := generateTemp()
	display := make(chan int)
	outputTemp(display)
	t := <-temps
	for {
		select {
		/*
		   since the temps gets sent new messages every 200 milliseconds and we only consume
		   a latest temp every 2 seconds, we consume the temps and set it to the t variable.
		   this will ensure that when the outputTemp goroutine is ready to consume the latest temp
		   and print it, it will output the latest value.
		*/
		// consuming temps and setting latest value to t here
		case t = <-temps:
			/*
			   when the display channel is attempting to consume a message (every 2 seconds), send
			   the last t value as that's the latest temp
			*/
		case display <- t:
		}
	}
}

func generateTemp() chan int {
	output := make(chan int)
	go func() {
		temp := 50
		for {
			output <- temp
			temp += rand.Intn(3) - 1
			time.Sleep(200 * time.Millisecond)
		}
	}()
	return output
}

func outputTemp(input chan int) {
	go func() {
		for {
			fmt.Println("Current temp:", <-input)
			time.Sleep(2 * time.Second)
		}
	}()
}
