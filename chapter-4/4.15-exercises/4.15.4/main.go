package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

func main() {
	mutex := sync.Mutex{}
	var wordFrequency = make(map[string]int)
	for i := 1000; i <= 1030; i++ {
		url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
		go countWordFrequency(wordFrequency, url, &mutex)
	}

	time.Sleep(time.Second * 10)
	mutex.Lock()
	for word, count := range wordFrequency {
		log.Printf("word: %v -> %d\n", word, count)
	}
	mutex.Unlock()
}

func countWordFrequency(wordFrequency map[string]int, url string, mutex *sync.Mutex) {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		panic("Server returning error status code: " + resp.Status)
	}

	body, _ := io.ReadAll(resp.Body)
	wordRegex := regexp.MustCompile(`[a-zA-Z]+`)

	words := wordRegex.FindAllString(string(body), -1)
	mutex.Lock()
	fmt.Printf("Updating word frequency for url: %v. Locking shared memory\n", url)
	for _, w := range words {
		w = strings.ToLower(w)
		wordFrequency[w] += 1
	}
	fmt.Printf("Finish updating word frequency for url: %v. Unlocking shared memory\n", url)
	mutex.Unlock()
}
