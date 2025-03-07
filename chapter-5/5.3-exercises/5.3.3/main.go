package main

import (
	"fmt"
	"time"

	"5.3.3/wsemaphore"
)

func main() {
	sema := wsemaphore.NewWeightedSemaphore(3)
	sema.Acquire(2)
	fmt.Println("Parent thread acquired semaphore")
	go func() {
		sema.Acquire(2)
		fmt.Println("Child thread acquired semaphore")
		sema.Release(2)
		fmt.Println("Child thread released semaphore")
	}()
	time.Sleep(3 * time.Second)
	fmt.Println("Parent thread releasing semaphore")
	sema.Release(2)
	time.Sleep(1 * time.Second)
}
