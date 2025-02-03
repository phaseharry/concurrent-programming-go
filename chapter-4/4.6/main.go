package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

const AllLetters = "abcdefghijklmnopqrstuvwxyz"

func main() {
	mutex := sync.Mutex{}
	var frequency = make([]int, 26)

	for i := 2000; i <= 2200; i++ {
		url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
		go countLettersConcurrent(url, frequency, &mutex)
	}

	// tries to log output 100 times
	for i := 0; i < 100; i++ {
		time.Sleep(100 * time.Millisecond) // sleep for 100 milliseconds
		/*
		 non-blocking. (main process will continue to do this as long as mutex is available).
		 It won't wait until mutex become available so it won't block other more important goroutines
		 that actually need the mutex to do processing with. It will only acquire the lock if it's available
		 and that's it.
		*/
		if mutex.TryLock() {
			for i, c := range AllLetters {
				fmt.Printf("%c-%d ", c, frequency[i])
			}
			mutex.Unlock()
		} else {
			// if mutex is not availabe, log it and try again (max 100 tries with a 100 millisecond delay for the main goroutine)
			fmt.Println("Mutex already inuse")
		}
	}
}

func countLettersConcurrent(url string, frequency []int, mutex *sync.Mutex) {
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		panic("Server returning error status code: " + resp.Status)
	}

	body, _ := io.ReadAll(resp.Body)

	/*
		There is a cost to calling Lock/Unlock since the software needs to communicate with the actual hardware
		to get exclusive access to the critical sections, so instead of calling Lock/Unlock for each character
		that gets incremented into our shared fequency slice, we can lock it for the entire file and read and increment
		all of its characters while we have exclusive access to processing with it.
		ex)

		mutex.Lock()
		if cIndex >= 0 {
			frequency[cIndex] += 1
		}
		mutex.Unlock()
	*/
	mutex.Lock()
	for _, b := range body {
		char := strings.ToLower(string(b))
		cIndex := strings.Index(AllLetters, char)
		if cIndex >= 0 {
			frequency[cIndex] += 1
		}
	}
	mutex.Unlock()
	fmt.Println("Completed: ", url, time.Now().Format("15:04:05"))
}
