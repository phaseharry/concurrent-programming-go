package main

import (
	"fmt"
	"sync"
	"time"
)

/*
Simple example showing a possbile deadlock happening.
2 goroutines compete for the same exclusive access to 2 mutexes.
Both goroutines A and B wants exclusive access to both lock 1 and 2.
If A was able to get lock1 and B was able to get lock2, then the application
will reach a state in which both A and B will be blocked and will not be able to continue,
leading to a deadlock.
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

func red(lock1, lock2 *sync.Mutex) {
	for {
		fmt.Println("Red: Acquiring lock1")
		lock1.Lock()
		fmt.Println("Red: Acquiring lock2")
		lock2.Lock()
		fmt.Println("Red: Both locks acquired")
		lock1.Unlock()
		lock2.Unlock()
		fmt.Println("Red: locks released")
	}
}

func blue(lock1, lock2 *sync.Mutex) {
	for {
		fmt.Println("Blue: Acquiring lock2")
		lock2.Lock()
		fmt.Println("Blue: Acquiring lock1")
		lock1.Lock()
		fmt.Println("Blue: Both locks acquired")
		lock1.Unlock()
		lock2.Unlock()
		fmt.Println("Blue: locks released")
	}
}
