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
demonstration of closing a channel after a condition has been met so
to as short circuit the pipeline.
*/
func main() {
	startTime := time.Now()
	/*
	   creating 2 quit signal channels.
	   - (quitWords) one is specifically used to quitWords once we
	   have extracted 10,000 words from our page contents. it will be used for any stages
	   of the pipeline up and until extractWords.
	   - (quit) is used for any of the stages after it once we've reached 10,000 words.

	   this will ensure that we will stop consuming words once we've hit 10,000 words but
	   any other goroutines that are processing those 10,000 words after the extractWords stage
	   will not be interrupted and continue processing those first 10,000 words.
	*/
	quitWords := make(chan int)
	quit := make(chan int)
	defer close(quit)
	urls := generateUrls(quitWords)
	pages := make([]<-chan string, DOWNLOADERS)
	for i := range DOWNLOADERS {
		pages[i] = downloadPages(quitWords, urls)
	}
	fannedInChannel := FanIn(quitWords, pages...)
	words := Take(
		quitWords,
		10_000,
		extractWords(quitWords, fannedInChannel),
	)
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

func Broadcast[K any](quit <-chan int, input <-chan K, n int) []chan K {
	outputs := CreateAll[K](n)

	go func() {
		defer CloseAll(outputs...)
		var message K
		moreData := true
		for moreData {
			select {
			case message, moreData = <-input:
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

func CreateAll[K any](n int) []chan K {
	channels := make([]chan K, n)
	for i, _ := range channels {
		channels[i] = make(chan K)
	}
	return channels
}

func CloseAll[K any](channels ...chan K) {
	for _, channel := range channels {
		close(channel)
	}
}

/*
takes as an input channel and adds a condition modifier to it.
it pipes the messages from input into a newly created output channel.
it takes a value n indicating the number of messages it will consume
before it short circuits and stop listening to messages from input
*/
func Take[K any](quit chan int, n int, input <-chan K) <-chan K {
	output := make(chan K)

	go func() {
		initialLimit := n
		defer close(output)
		moreData := true
		var message K
		for n > 0 && moreData {
			select {
			case message, moreData = <-input:
				if moreData {
					output <- message
					n -= 1
				}
			case <-quit:
				return
			}
		}
		if n == 0 {
			fmt.Printf(
				"reached take limit of %v, closing quit channel to short circuit\n",
				initialLimit,
			)
			close(quit)
		}
	}()

	return output
}
