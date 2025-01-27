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

	for i := 1000; i <= 1030; i++ {
		url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
		// go countLettersSequential(url, frequency, &mutex)
		go countLettersConcurrent(url, frequency, &mutex)
	}

	time.Sleep(time.Second * 60)
	mutex.Lock()
	for i, c := range AllLetters {
		fmt.Printf("%c-%d\n", c, frequency[i])
	}
	mutex.Unlock()
}

// incorrect way of using mutex
func countLettersSequential(url string, frequency []int, mutex *sync.Mutex) {
	/*
		getting exclusive access here when we're not reading/writing to the shared memory resource.
		This will block other goroutines that are just fetching for the file first and not updating the shared
		memory resource, effectively making our program into a sequential program instead of a concurrent one
		since all operations has to happen one after another. We can solve this by only locking the parts
		where we're actually reading/writing to shared memory resource so the fetching of the file can still
		be done concurrently/in parallel.
	*/
	mutex.Lock()
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		panic("Server returning error status code: " + resp.Status)
	}

	body, _ := io.ReadAll(resp.Body)
	for _, b := range body {
		char := strings.ToLower(string(b))
		cIndex := strings.Index(AllLetters, char)
		if cIndex >= 0 {
			frequency[cIndex] += 1
		}
	}

	fmt.Println("Completed: ", url, time.Now().Format("15:04:05"))
	mutex.Unlock()
}

// correct way
func countLettersConcurrent(url string, frequency []int, mutex *sync.Mutex) {
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		panic("Server returning error status code: " + resp.Status)
	}

	body, _ := io.ReadAll(resp.Body)

	/*
		There is a cost to calling Lock/Unlock since the software needs to communicate with the actual hardware
		to get exclusive access to the critical sectionsm, so instead of calling Lock/Unlock for each character
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
