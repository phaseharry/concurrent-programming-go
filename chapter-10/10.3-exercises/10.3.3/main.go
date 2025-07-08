package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

func main() {
	const pagesToDownload = 30
	lineResults := make(chan int)
	wg := &sync.WaitGroup{}
	wg.Add(pagesToDownload)

	for i := 1000; i < 1000+pagesToDownload; i++ {
		go func(id int) {
			url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", id)
			fmt.Println("Downloading", url)
			resp, _ := http.Get(url)
			if resp.StatusCode != 200 {
				panic("Server’s error: " + resp.Status)
			}
			bodyBytes, _ := io.ReadAll(resp.Body)
			lineResults <- strings.Count(string(bodyBytes), "\n")
			resp.Body.Close()
			wg.Done()
		}(i)
	}

	result := joinResults(lineResults)

	wg.Wait()

	/*
	   closing lineResults once we're done processing urls and queued their line numbers
	   into the lineResults channel. Closing so in joinResults, we'll be unlucked in the
	   goroutine created there so it can return the total lines sum
	*/
	close(lineResults)

	fmt.Println("Total lines:", <-result)
}

func joinResults(lineResults chan int) chan int {
	result := make(chan int)

	go func() {
		totalLines := 0
		for lines := range lineResults {
			totalLines += lines
		}
		result <- totalLines
	}()

	return result
}
