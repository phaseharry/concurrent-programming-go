package main

import (
	"fmt"

	"5.11/semaphore"
)

/*
Semaphores can reduce the likely hood of deadlocks due to goroutines
that are waiting missing the Signal or Broadcast message (see ex. 5.7).
The ordering of the Acquire / Release doesn't matter due to the permit being
an int determining whether to block a goroutine's execution or not.

ex. If the doWork goroutine is ran and calls .Release() before Acquire() is called,
then the permit counter would be 1. When Acquire() is called after, the permit counter
will be 1 so it can acquire access to the semaphore and execute.

The same would be happen if Acquire() is called before Release().
If Acquire() is called first, the permit would be 0, so it would called .Wait()
and be blocked. Then Release() is called and increment the permit to 1 and call
.Signal() internally so the goroutine with Acquire() called would pick it up.
The goroutine that called .Acquire() would check the permit again and see
that it is above 0 and can execute. Internally, the .Acquire() method would
also decrement the permit value to keep track of outstanding goroutine executions
using the semaphore.
*/
func main() {
	semaphore := semaphore.NewSemaphore(0)

	for range 50000 {
		go doWork(semaphore)
		fmt.Println("Waiting for child goroutine")
		semaphore.Acquire()
		fmt.Println("Child goroutine finished")
	}
}

func doWork(semaphore *semaphore.Semaphore) {
	fmt.Println("Work started")
	fmt.Println("Work finished")
	semaphore.Release()
}
