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
	wg := sync.WaitGroup{}
	wg.Add(31)
	mutex := sync.Mutex{}
	var frequency = make([]int, 26)
	for i := 1000; i <= 1030; i++ {
		url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
		/*
		   creating an anonymous function to handle the goroutine of calling
		   the api and incrementing the frequencies for its text and calling
		   the .Done() call for our wait group. this ensures we don't have to
		   modify the countLettersConcurrent function we created in chapter 4.
		*/
		go func() {
			countLettersConcurrent(url, frequency, &mutex)
			fmt.Println("finished incrementing frequency for link:", i)
			wg.Done()
		}()
	}
	/*
	   wait until all 31 api calls to fetch for text and
	   increment the character frequency is done
	*/
	wg.Wait()

	mutex.Lock()
	// printing the results of the characters from all 31 texts
	for i, c := range AllLetters {
		fmt.Printf("%c-%d ", c, frequency[i])
	}
	mutex.Unlock()
}

func countLettersConcurrent(url string, frequency []int, mutex *sync.Mutex) {
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		panic("Server returning error status code: " + resp.Status)
	}

	body, _ := io.ReadAll(resp.Body)
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
