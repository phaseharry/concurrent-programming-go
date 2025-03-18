package main

import (
	"fmt"
	"time"

	"6.3.2/waitgrp"
)

func main() {
	wg := waitgrp.NewWaitGrp()
	wg.Add(1)

	/*
	   spin off a child goroutine that sleeps for 5 seconds and
	   then calls .Done() on WaitGroup and lets the main goroutine
	   terminate
	*/
	go func() {
		fmt.Println("Sleeping for 5 seconds")
		time.Sleep(5 * time.Second)
		fmt.Println("Marking waitgroup as done")
		wg.Done()
	}()
	/*
	   using the TryWait implemented in WaitGroup. TryWait will
	   return a bool indicating whether the WaitGroup counter has hit 0 or not.
	   If WaitGroup is still waiting, then have main goroutine sleep for a second
	   and try again after.
	*/
	for !wg.TryWait() {
		fmt.Println("Wait group is done yet. Try again after 1 second delay")
		time.Sleep(time.Second)
	}

	fmt.Println("Wait group is done")
}
