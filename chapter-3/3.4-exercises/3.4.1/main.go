package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func main() {
	var wordFrequency = make(map[string]int)
	for i := 1000; i <= 1030; i++ {
		url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
		go countWordFrequency(wordFrequency, url)
	}
	time.Sleep(time.Second * 10)
	for word, count := range wordFrequency {
		log.Printf("word: %v -> %d\n", word, count)
	}
}

func countWordFrequency(wordFrequency map[string]int, url string) {
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
	for _, w := range words {
		w = strings.ToLower(w)
		wordFrequency[w] += 1
	}
}
