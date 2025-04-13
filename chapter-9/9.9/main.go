package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

const DOWNLOADERS = 20

func main() {
	startTime := time.Now()
	quit := make(chan int)
	defer close(quit)
	urls := generateUrls(quit)

	/*
	   create a slice of capacity 20 that holds the 20 channels that will be used
	   by 20 downloader goroutines to fetch for page content concurrently. This is
	   demonstrating the "Fanout" pattern in which one goroutine's results gets
	   split up into multiple goroutine in the next step of the pipeline.

	   ex. having urls be created and sent through one channel (urls) and then
	   spliting the results of the generated urls into 20 seperate goroutines to be processed concurrently.

	   Fetching for content through the internet is a time consuming process, so doing it concurrently
	   is ideal and will speed up the application
	*/
	pages := make([]<-chan string, DOWNLOADERS)
	for i := range DOWNLOADERS {
		pages[i] = downloadPages(quit, urls)
	}

	/*
	   since we have 20 goroutines fetching for contents and outputting it through 20 seperate channels,
	   we need a way to merge the results from the 20 seperate channels into one channel to pipe it to the
	   extractWords goroutine in the pipeline. there are multiple options for doing this.

	   - We could've created a single channel that is passed into to each downloadPages goroutine. this common
	   channel would be used as the output channel that each one of the downloadPages send their fetched pageContent
	   into. this would've fanned in all of the 20 goroutines results into that channel and that channel would be piped into
	   the extraWords goroutine.

	   - actual implementation: keep similar pattern as current pipeline implementation which involves passing a quit channel
	   and the initial input channel and having it return the output channel. Each of the 20 downloadPages goroutine will
	   receive a quit channel and the input urls channel. each goroutine will consume a url has it becomes available. once a
	   goroutine consumes a url, it will not be available to be consumed by another goroutine. to merge in all the pageContents
	   received from the 20 different downloadPages goroutines, we have a fanIn helper function that will consume messsages from
	   all 20 channels for pageContent and pipe it into one common channel to be sent to the extractWords goroutine for further processing.
	*/
	fannedInChannel := fanIn(quit, pages...)

	results := extractWords(quit, fannedInChannel)
	for result := range results {
		fmt.Println(result)
	}
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

func fanIn[K any](quit <-chan int, allChannels ...<-chan K) chan K {
	wg := sync.WaitGroup{}
	/*
	   since we are depending on 20 goroutines that spawn 20 channels, we can't simply listen to only
	   one channel to see if it's active or not, we need to listen to all 20 channels from downloadPages.
	   if even one is active, we have to keep consuming pageContent, so to decide when to terminate / know
	   we're done with the fanIn goroutine, we will use a WaitGroup that's waiting on 20 waitGroup.Done()
	   calls to know when the fanIn goroutine is done processing.
	*/
	wg.Add(len(allChannels))

	// output channel that gets messaged fanned in to
	output := make(chan K)
	for _, c := range allChannels {
		/*
		   for each channel that fans into our output channel, we start a seperate goroutine for each channel to consume
		   messages from the respective channel. We consume messages from that channel until it's not active anymore or the
		   quit signal was sent. once either happens, the respective goroutine will call Wait.Done() to signal to the waitGroup
		   that it's waiting on one less .Done() call.
		*/
		go func(channel <-chan K) {
			defer wg.Done()
			for i := range channel {
				select {
				case output <- i:
				case <-quit:
					return
				}
			}
		}(c) // passing in each channel to the anonymous function so code knows what to reference
	}
	/*
	   start another goroutine that is blocked until all the channels have their messages consumed and are closed.
	   Once all channels have been closed, we will close the output channel indicating that there is no more messages
	   to be sent to it.

	   we start another goroutine and have it block to prevent an earliest closure of the output channel because we want to
	   send the output channel out to be consumed and have its messages be processed in the main goroutine.
	*/
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
