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
	for i := 0; i < 20; i++ {
		cond.L.Lock()

		for *money < 50 {
			/*
				if our shared "money" variable is less than 50, then we have our condition
				call Wait(). The Wait() call unlocks the mutex so another goroutine that is waiting
				for exclusive access to the share resource can process and update the shared resource.
				In this case, it is the stingy function. After calling .Wait(), this goroutine will be
				waiting for a .Signal() call from the other goroutine. That call will signify that the
				other goroutine has made a change to the shared resource (updated its value in some way)
				that might pass the "condition" you have set for this goroutine to be able to process
				(in this case, it is that there needs to be at least 50 in money to be able to subtract 50
				from the money variable.)
				Once that resource calls .Signal and unlocks the mutex, this goroutine will be able to acquire the
				log and attempt another pass.
				if we got a .Signal call, but it is still not at least 50, we will call .Wait() again and the
				above repeats until we pass the condition.
			*/
			cond.Wait()
		}

		*money -= 50

		// case should never happen since we only subtract IF we have at least 50 so smallest
		// value money can be is 0
		if *money < 0 {
			fmt.Println("Money is negative!!")
			os.Exit(1)
		}

		cond.L.Unlock()

	}
	fmt.Println("Spendy Done")
}

func stingy(money *int, cond *sync.Cond) {
	for i := 0; i < 100; i++ {
		/*
			Acquire the lock within the condition. If we can't acquire the lock
			then this goroutine will be blocked and wait until it can acquire the lock.
			Once we add 10 to the shared "money" variable, we call the condition Signal method so any
			goroutine that has called Wait method on the condition struct can attempt to acquire
			the Lock again and reprocess based on some condition.
		*/
		cond.L.Lock()
		*money += 10
		cond.Signal()

		// will have to call unlock so the signaled goroutine can actually acquire lock and attempt another attempt at processing
		cond.L.Unlock()
	}
	fmt.Println("Stingy Done")
}
