package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

func main() {
	money := 100
	mutex := sync.Mutex{}
	condition := sync.NewCond(&mutex)

	go spendy(&money, condition)
	go stingy(&money, condition)

	time.Sleep(2 * time.Second)
	mutex.Lock()
	fmt.Println("Money in bank account: ", money)
	mutex.Unlock()
}

func spendy(money *int, cond *sync.Cond) {
	for range 20 {
		cond.L.Lock()

		// have spendy goroutine execution wait until we have at least 50 dollars
		for *money < 50 {
			cond.Wait()
		}

		*money -= 50

		if *money < 0 {
			fmt.Println("Money is negative!!")
			os.Exit(1)
		}

		cond.L.Unlock()

	}
	fmt.Println("Spendy Done")
}

func stingy(money *int, cond *sync.Cond) {
	for range 100 {
		cond.L.Lock()
		*money += 10

		/*
			only signal if there's 50 dollars or more so the spendy goroutines don't arbitrarily
			attempt to get access and try to execute even if it'll just be blocked again.
			this will let the stingy goroutine continually have access and process.
		*/
		if *money >= 50 {
			cond.Signal()
		}

		cond.L.Unlock()
	}
	fmt.Println("Stingy Done")
}
