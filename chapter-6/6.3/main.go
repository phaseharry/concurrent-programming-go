package main

import (
	"fmt"

	"6.3/waitgrp"
)

func main() {
	/*
	   setting the internal semaphore permit count to -3
	   (1 - 4) = -3, so when each child goroutine calls .Done()
	   the permit count with increment by 1. Once we're at permit
	   count of 1, then the main goroutine will be unlocked and can finish.
	*/
	wg := waitgrp.NewWaitGrp(4)
	for i := 1; i <= 4; i++ {
		go doWork(i, wg)
	}
	/*
	   blocks main goroutine until the internal semaphore within
	   the waitGroup has a permit value of 1.
	*/
	wg.Wait()
	fmt.Println("All complete")
}

func doWork(id int, wg *waitgrp.WaitGrp) {
	fmt.Println(id, "Done working")
	wg.Done()
}
