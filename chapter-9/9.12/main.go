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
demonstation of flushing results when closed.
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
	topTenLongestWords := longestWords(quit, words)
	results := <-topTenLongestWords
	fmt.Println(results)
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
		/*
		   goroutine that is blocked entirely until we consumed all words from all pageContents
		   so we have all possible words in the pageContents to compare lengths to. this is done
		   so we can get an accurate length of the top 10 longest words.
		   once we're done and processed all words, we sort the words by length and return a slice of it
		   as output. we don't send results as we process a word. we only send when we've calculated everything,
		   i.e. "flushing on close"
		*/
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
