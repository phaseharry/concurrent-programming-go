package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	mutex := sync.Mutex{}
	count := 5
	go countdown(&count, &mutex)

	mutex.Lock()
	remainingCount := count
	mutex.Unlock()

	for remainingCount > 0 {
		time.Sleep(500 * time.Millisecond)
		mutex.Lock()
		fmt.Println("main goroutine:", remainingCount)
		remainingCount = count
		mutex.Unlock()
	}
}

func countdown(seconds *int, mutex *sync.Mutex) {
	mutex.Lock()
	remainingSeconds := *seconds
	mutex.Unlock()

	for remainingSeconds > 0 {
		time.Sleep(time.Second * 1)
		mutex.Lock()
		*seconds -= 1
		fmt.Println("countdown goroutine:", *seconds)
		remainingSeconds = *seconds
		mutex.Unlock()
	}
}
