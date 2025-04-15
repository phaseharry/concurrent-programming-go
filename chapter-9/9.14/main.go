package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

const DOWNLOADERS = 20

/*
demonstation of broadcasting one channels output into mutliple goroutines
so goroutines can process seperate jobs in parallel with the same available data
*/
func main() {
	startTime := time.Now()
	quit := make(chan int)
	defer close(quit)
	urls := generateUrls(quit)

	pages := make([]<-chan string, DOWNLOADERS)
	for i := range DOWNLOADERS {
		pages[i] = downloadPages(quit, urls)
	}

	fannedInChannel := FanIn(quit, pages...)

	words := extractWords(quit, fannedInChannel)
	/*
	   taking the single words channel and broadcasting its values into 2 seperate
	   channels to be consumed individually by its respective goroutines in the pipeline
	*/
	wordsBroadcasted := Broadcast(quit, words, 2)
	topTenLongestWords := longestWords(quit, wordsBroadcasted[0])
	topTenFrequentWords := frequentWords(quit, wordsBroadcasted[1])

	fmt.Println("Top 10 Longest Words:", <-topTenLongestWords)
	fmt.Println("Top 10 Most Frequent Words:", <-topTenFrequentWords)
	duration := time.Since(startTime)
	fmt.Println("concurrent page download duration:", duration)
}

func generateUrls(quit <-chan int) <-chan string {
	urls := make(chan string)

	go func() {
		defer close(urls)
		for i := 100; i <= 130; i++ {
			url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
			select {
			case urls <- url:
			case <-quit:
				return
			}
		}
	}()

	return urls
}

func downloadPages(quit <-chan int, urls <-chan string) <-chan string {
	pages := make(chan string)

	go func() {
		defer close(pages)
		urlsChannelActive, url := true, ""

		for urlsChannelActive {
			select {
			case url, urlsChannelActive = <-urls:
				if urlsChannelActive {
					resp, _ := http.Get(url)
					if resp.StatusCode != 200 {
						panic("Server's error: " + resp.Status)
					}
					body, _ := io.ReadAll(resp.Body)
					pages <- string(body)
					resp.Body.Close()
				}
			case <-quit:
				return
			}
		}
	}()

	return pages
}

func extractWords(quit <-chan int, pages <-chan string) <-chan string {
	words := make(chan string)

	go func() {
		defer close(words)
		wordRegex := regexp.MustCompile(`[a-zA-Z]+`)
		pageChannelActive, pageContent := true, ""

		for pageChannelActive {
			select {
			case pageContent, pageChannelActive = <-pages:
				if pageChannelActive {
					for _, word := range wordRegex.FindAllString(pageContent, -1) {
						words <- strings.ToLower(word)
					}
				}
			case <-quit:
				return
			}
		}
	}()

	return words
}

func longestWords(quit <-chan int, words <-chan string) <-chan string {
	results := make(chan string)

	go func() {
		defer close(results)
		uniqueWordsMap := make(map[string]bool)
		uniqueWords := make([]string, 0)
		moreData, word := true, ""
		for moreData {
			select {
			case word, moreData = <-words:
				if moreData && !uniqueWordsMap[word] {
					uniqueWordsMap[word] = true
					uniqueWords = append(uniqueWords, word)
				}
			case <-quit:
				return
			}
		}

		sort.Slice(uniqueWords, func(a, b int) bool {
			return len(uniqueWords[a]) > len(uniqueWords[b])
		})

		results <- strings.Join(uniqueWords[:10], ", ")
	}()

	return results
}

func frequentWords(quit <-chan int, words <-chan string) <-chan string {
	mostFrequentWords := make(chan string)

	go func() {
		defer close(mostFrequentWords)
		freqMap := make(map[string]int)
		freqList := make([]string, 0)
		moreData, word := true, ""
		for moreData {
			select {
			case word, moreData = <-words:
				if moreData {
					if freqMap[word] == 0 {
						freqList = append(freqList, word)
					}
					freqMap[word] += 1
				}
			case <-quit:
				return
			}
		}
		sort.Slice(freqList, func(a, b int) bool {
			return freqMap[freqList[a]] > freqMap[freqList[b]]
		})
		mostFrequentWords <- strings.Join(freqList[:10], ", ")
	}()

	return mostFrequentWords
}

func FanIn[K any](quit <-chan int, allChannels ...<-chan K) chan K {
	wg := sync.WaitGroup{}
	wg.Add(len(allChannels))

	output := make(chan K)
	for _, c := range allChannels {
		go func(channel <-chan K) {
			defer wg.Done()
			for i := range channel {
				select {
				case output <- i:
				case <-quit:
					return
				}
			}
		}(c)
	}
	go func() {
		wg.Wait()
		close(output)
	}()

	return output
}

/*
takes one channel as input and creates n channels to pipe the messages
consumed from the input channel into the newly created n channels.
these channels will be sent to seperate goroutines to be consumed individually
*/
func Broadcast[K any](quit <-chan int, input <-chan K, n int) []chan K {
	outputs := CreateAll[K](n)

	go func() {
		defer CloseAll(outputs...)
		var message K
		moreData := true
		for moreData {
			select {
			case message, moreData = <-input:
				/*
					consumes a message from the input channel and iterate through our list of output channels
					and send that message to each of those channels
				*/
				if moreData {
					for _, broadcastChannel := range outputs {
						broadcastChannel <- message
					}
				}
			case <-quit:
				return
			}
		}
	}()

	return outputs
}

/*
utility to create n goroutines of type K.
using Generics so we can reuse CreateAll for creating
many channels of the same type
*/
func CreateAll[K any](n int) []chan K {
	channels := make([]chan K, n)
	for i, _ := range channels {
		channels[i] = make(chan K)
	}
	return channels
}

// utility to close multple channels that comes as a slice
func CloseAll[K any](channels ...chan K) {
	for _, channel := range channels {
		close(channel)
	}
}
