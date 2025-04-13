package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

/*
demonstration of pipeling messages from goroutine to goroutine using channels.
each goroutine in the chain is responsible for a specific task.
once that task is done, then it closes the channel that it is responsible for.
idea is each goroutine consumes a channel created within the previous channel in the chain,
processes the messages it receives, and send the results of that processed image into
the next goroutine in the chain to further process

each goroutine always gets passed a quit channel which gets sent no messages.
We will use it as a signaling channel to abort / stop existing goroutines in the
pipeline so stop processing for whatever reason. we do this by closing the quit channel
at which time consumption of messages from the quit channel will start receiving the
default int value of 0 and we stop processing
*/
func main() {
	startTime := time.Now()
	quit := make(chan int)
	defer close(quit)
	urls := generateUrls(quit)
	pages := downloadPages(quit, urls)
	words := extractWords(quit, pages)
	for word := range words {
		fmt.Println(word)
	}
	duration := time.Since(startTime)
	fmt.Println("sequential page download duration:", duration)
}

func generateUrls(quit <-chan int) <-chan string {
	urls := make(chan string)

	go func() {
		/*
		   once we're done sending 30 urls through the urls channel,
		   we will close the urls channel so any goroutines that consume
		   urls can be signaled that we're done with the urls and stop
		   consuming. they can even stop processing if its done with its
		   existing tasks
		*/
		defer close(urls)
		for i := 100; i <= 130; i++ {
			url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
			select {
			/*
			   initial goroutine in the chain, it generates 30 urls of content we want to fetch
			   and extract the words out of. We will generate the url and send it through the url channel
			   so the next goroutine in the pipeline can consume it.
			*/
			case urls <- url:
				/*
				   if we get a quit message then we know the quit signal was set (channel was closed)
				   so we stop processing immediately and return
				*/
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
		/*
		   same as generateUrls goroutine, close the pages channel when done fetching pageContent
		   so the next goroutine in the pipeline will know when it's done processing and terminate
		*/
		defer close(pages)
		urlsChannelActive, url := true, ""

		/*
		   the second goroutine in the pipeline, it consumes urls from the urls channel
		   and fetches for page content. Like the generateUrls goroutine, it also gets
		   the quit channel that can preemptively signal to stop processing if an issue
		   occurred in another goroutine.

		   we only continue when the urls channel is active, indicating that there's more
		   urls to fetch the page content of
		*/
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
		/*
		   closes words channel when done processing so it can signal to the main goroutine that
		   we have sent all words and it can terminate the program once all words are consumed
		*/
		defer close(words)
		// regex to match every word entry within a page content
		wordRegex := regexp.MustCompile(`[a-zA-Z]+`)
		pageChannelActive, pageContent := true, ""
		/*
		   3rd goroutine in pipeline, consumes the pageContents from the pages channel and extract the words
		   from it. Send each word through the words channel.
		   receives the quit channel so it can be signaled to preemptively terminate if necessary
		*/
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
