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
	go spendy(&money, &mutex)
	go stingy(&money, &mutex)

	time.Sleep(2 * time.Second)
	mutex.Lock()
	fmt.Println("Money in bank account: ", money)
	mutex.Unlock()
}

func spendy(money *int, mutex *sync.Mutex) {
	/*
		does 1/5 of the iterations of stingy function, but money
		should still be 0 since each iteration of spendy is -50
		while each iteration of stingy is +10. However, we might reach
		negative territory since there's 2 goroutines sharing the share memory
		resource of "money". We will most likely reach negative territory
		with our money value which we don't want, so we exit with an OS call.
		See example 5.2 for an example of blocking and waiting if our money is at
		0. In that example we will wait for stingy to add enough money so stingy
		function can spend again.
	*/
	for i := 0; i < 20; i++ {
		mutex.Lock()
		*money -= 50
		if *money < 0 {
			fmt.Println("Money is negative!")
			os.Exit(1)
		}
		mutex.Unlock()
	}
	fmt.Println("Spendy Done")
}

func stingy(money *int, mutex *sync.Mutex) {
	for i := 0; i < 100; i++ {
		mutex.Lock()
		*money += 10
		mutex.Unlock()
	}
	fmt.Println("Stingy Done")
}
