package main

import (
	"fmt"
	"sync"
)

/*
this is an example of a deadlock in which one of the doWork
goroutine is not able to get the lock because there was a .Wait()
from the goroutine so it's waiting for another goroutine to .Signal()
or .Boardcast(). However, in this scenario, the .Signal call was missed
(it was called when there was no goroutine waiting for it) so when the
next goroutine calls .Wait() it is waiting forever and a .Signal call will
not happen.

To ensure that this does not happen, when .Signal() or .Broadcast() is called,
there needs to be another goroutine waiting for it; otherwise the signal or broard
is not received by any goroutine and it's missed.
*/
func main() {
	cond := sync.NewCond(&sync.Mutex{})
	cond.L.Lock()

	// modernized version of iterating through something 5000 times if
	// no idx value is used for something
	for range 50000 {
		/*
			issue occurs due to the child goroutine doWork calls Signal()
			in this manner. If the main goroutine's Wait() ever happens AFTER
			the Signal() call then the main goroutine will be stuck in a waiting
			state and not proceed any further since it "missed" the signal to acquire
			the lock and do the next iteration.
		*/
		go doWork(cond)
		fmt.Println("Waiting for child goroutine")
		cond.Wait()
		fmt.Println("Child goroutine finished")
	}
	cond.L.Unlock()
}

func doWork(cond *sync.Cond) {
	fmt.Println("Work started")
	fmt.Println("Work finished")
	cond.Signal()
}
