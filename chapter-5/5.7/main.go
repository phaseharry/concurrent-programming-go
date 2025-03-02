package main

import (
	"fmt"
	"sync"
)

/*
Solution to the issues we've seen in 5.6 in which
the main goroutine missed the doWork child goroutine's
.Signal() call, leading to a deadlock.
*/
func main() {
	cond := sync.NewCond(&sync.Mutex{})
	cond.L.Lock()
	for range 50000 {
		go doWork(cond)
		fmt.Println("Waiting for child goroutine")
		cond.Wait()
		fmt.Println("Child goroutine finished")
	}
}

/*
To prevent the main goroutine from calling .Wait()
before the child doWork goroutine calls .Signal(),
we ensure we have exclusive access to the mutex.
Since the main goroutine still has the lock from line 15,
the child goroutine will be blocked until the mutex is unlocked.
This will let the main goroutine call .Wait() before this doWork
goroutine from calling .Signal since .Wait() will release the
lock so the doWork goroutine can acquire the lock and call .Signal().
Once that's called and it also Unlocks the mutual exclusive access,
the main goroutine's .Wait() will acquire the lock again and iterate to
the next call.
*/
func doWork(cond *sync.Cond) {
	fmt.Println("Work started")
	fmt.Println("Work finished")
	cond.L.Lock()
	cond.Signal()
	cond.L.Unlock()
}
