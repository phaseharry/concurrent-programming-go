package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	money := 100
	/*
		Creating a new mutex that will protect the critical sections when
		we update the money variable concurrrently through multiple goroutines.
		This is done to prevent race conditions that might overwrite existing values incorrectly.
		The goroutines will use this mutex to lock the critical section in question
		down so when a goroutine has a lock on it, no other goroutines can run that section
		until an unlock is called to unblock the execution for other goroutines.
	*/
	mutex := sync.Mutex{}

	go stingy(&money, &mutex)
	go spendy(&money, &mutex)

	/*
		manually sleeping or the program will terminate when the main goroutine is done running even when
		there's still existing goroutines processing
	*/
	time.Sleep(time.Second * 2)

	/*
	 even if we're just reading the variable & not updating it, we still need to get a lock on the
	 mutex that protects it since the compiler might do some optimization during compile time
	 and the ordering for the print statement can happen before our execution fo stingy & spendy is
	 done running and give us the wrong answer
	*/
	mutex.Lock()
	fmt.Println("Money in the bank account: ", money)
	mutex.Unlock()
}

func stingy(money *int, mutex *sync.Mutex) {
	for i := 0; i < 1000000; i++ {
		mutex.Lock()
		*money += 10
		mutex.Unlock()
	}
	fmt.Println("Stingy done")
}

func spendy(money *int, mutex *sync.Mutex) {
	for i := 0; i < 1000000; i++ {
		mutex.Lock()
		*money -= 10
		mutex.Unlock()
	}
	fmt.Println("Spendy done")
}
