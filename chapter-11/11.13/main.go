package main

import (
	"fmt"
	"sync"
	"time"
)

/*
if you know the exact exclusive resources that concurrent executions will use in your program,
you can prevent deadlocks by having different concurrent executions acquire the exclusive resources
in the same order so that no goroutines can have only 1 of the required exclusive resource needed
to execute while the other have other required exclusive resources locked, preventing both goroutines
from executing further.
*/

func main() {
	lockA := sync.Mutex{}
	lockB := sync.Mutex{}
	go red(&lockA, &lockB)
	go blue(&lockA, &lockB)
	// not blocking the main goroutine so it will eventually terminate after 20 seconds
	time.Sleep(20 * time.Second)
	fmt.Println("Done")
}

/*
both red and blue goroutines acquire the shared exclusive resources in the same order to prevent
one from getting one lock while the other having the other lock. this will prevent a deadlock.

red -> gets lock1
blue -> attempts lock1, but can't acquire it so wait
red -> gets lock2
red -> finishes and releases both.
blue -> able to acquire lock1 so acquire it
red -> attempts to acquire lock1 again, but can't so blocked until lock1 is freed
blue -> acquires lock2
blue -> releases both lock1 and lock2
*/
func red(lock1, lock2 *sync.Mutex) {
	for {
		fmt.Println("Red: Acquiring lock1")
		lock1.Lock()
		fmt.Println("Red: Acquiring lock2")
		lock2.Lock()
		fmt.Println("Red: Both locks acquired")
		lock1.Unlock()
		lock2.Unlock()
		fmt.Println("Red: Locks released")
	}
}

func blue(lock1, lock2 *sync.Mutex) {
	for {
		fmt.Println("Blue: Acquiring lock1")
		lock1.Lock()
		fmt.Println("Blue: Acquiring lock2")
		lock2.Lock()
		fmt.Println("Blue: Both locks acquired")
		lock1.Unlock()
		lock2.Unlock()
		fmt.Println("Blue: Locks released")
	}
}
