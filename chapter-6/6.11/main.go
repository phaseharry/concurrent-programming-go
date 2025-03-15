package main

import (
	"fmt"
	"time"

	"6.11/barrier"
)

func main() {
	b := barrier.NewBarrier(2)
	/*
	   have 2 goroutines required to hit Barrier size and be allowed to continue processing.
	   Since the "Red" goroutine has a work duration of 4 seconds, it will always call .Wait()
	   first and be blocked until the "Blue" goroutine finishes its work which takes 10 seconds.
	   Thus "Red" will usually have to wait 6 seconds before "Blue" calls .Wait() and reaching
	   the Barrier waitSize and have both goroutines be allowed to do another round of iteration
	   within their infinite loops
	*/
	go workAndWait("Red", 4, b)
	go workAndWait("Blue", 10, b)

	/*
		letting the main goroutine run for 100 seconds so
		Read and Blue functions and work and wait
	*/
	time.Sleep(100 * time.Second)

}

func workAndWait(name string, timeToWork int, barrier *barrier.Barrier) {
	start := time.Now()
	/*
	   goroutine to simulate work being done.
	   prints interval between start and end where
	   there's a sleep in between. then calls .Wait()
	   to block until other goroutines have finished their
	   work and calls .Wait() to reach Barrier size.
	*/
	for {
		fmt.Println(time.Since(start), name, "is running")
		time.Sleep(time.Duration(timeToWork) * time.Second)
		fmt.Println(time.Since(start), name, "is waiting on barrier")
		barrier.Wait()
	}
}
