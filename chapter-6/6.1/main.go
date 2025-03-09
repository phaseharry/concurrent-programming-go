package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(4)
	for i := 1; i <= 4; i++ {
		go doWork(i, &wg)
	}
	/*
	   have the main goroutine block and wait until all child
	   goroutines (4) calls the .Done() of the waitGroup that's
	   synchronizing the operations.
	*/
	wg.Wait()
	fmt.Println("All complete")
}

func doWork(id int, wg *sync.WaitGroup) {
	/*
	   having each child goroutine sleep for a random interval
	   between 0 and 5 seconds and then calling .Done()
	*/
	i := rand.Intn(5)
	time.Sleep(time.Duration(i) * time.Second)
	fmt.Println(id, "Done working after", i, "seconds")
	wg.Done()
}
