package main

import (
	"fmt"
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
	for i := 0; i < 20; i++ {
		mutex.Lock()

		/*
			Instead of subtracting 50 and stopping the program when we have any
			dollars at negative, we will check if we have a dollar amount that's
			greater or equal to 50.
			- If we do, then we skip the for loop below and just subtract 50 from shared money variable,
			unlock and let the stingy goroutine have access to the shared resource.
			- if we don't have at least 50 within the money variable, then we
			release the lock and time this goroutine out for 10 milliseconds. This
			will let other goroutines (stingy) have access to the shared resource and
			add 10 dollars to it. Once 10 milliseconds is up and if we're able to
			reacquire the lock, we check if money is at least 50. If not, we continue
			the above logic. If we do, then we -50 and then release the lock.

			This solution would waste time as the 10 milliseconds might be too long as
			the other goroutine might've updated money and got at least 50 1 millisecond ago.

			What if we removed the sleep invocation? See below.
		*/
		for *money < 50 {
			mutex.Unlock()
			time.Sleep(10 * time.Millisecond)
			mutex.Lock()
		}

		/*
			below is the same as above, except there is no sleep time. Instead we just
			release the lock so another goroutine can acquire the lock. It could easily
			be acquire by the same goroutine that let it go even if the money shared resource
			was not updated at all.

			This solution would waste CPU resources since the CPU would be cycling needlessly
			and nothing would be changed since the money variable has not been updated by the
			stingy goroutine.

			See 5.3 for a better way to do this instead of adding an arbitrary sleep invocation
			or just iterating non-stop until the money variable has at least 50.

			for *money < 50 {
				mutex.Unlock()
				mutex.Lock()
			}
		*/

		*money -= 50

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
