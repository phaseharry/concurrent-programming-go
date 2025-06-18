package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

/*
using fork and join concurrency pattern to handle getting the code depth
of files within a directory, but having the main goroutine block until all
files have their code depths calculated and then getting the max code depth within
that directory
*/

type CodeDepth struct {
	file  string
	level int
}

func main() {
	directory := os.Args[1]
	codeDepths := make(chan CodeDepth)
	wg := sync.WaitGroup{}
	/*
		iterates through all files within directory and calls the passed
		in anonoymous function and forks to a different goroutine if it's a file.
		if it's subdirectory, ignore it.
	*/
	filepath.Walk(
		directory,
		func(path string, info os.FileInfo, err error) error {
			forkIfNeeded(path, info, &wg, codeDepths)
			return nil
		},
	)

	resultChannel := joinResults(codeDepths)

	wg.Wait()

	close(codeDepths)
	/*
		closing the codeDepths channel once all files have had their code depths
		calculated so the joinResults goroutine that is consuming codeDepth entries
		from the codeDepth channel will stop and send the max value over to the main goroutine
		in the line below.

		if we don't close then that loop in joinResults will not end and return the max value.
		that loop will be stuck there waiting to consume more codeDepth entries and be blocked,
		while the main goroutine will be blocked, waiting to consume the max depth, leading to a deadlock.
	*/
	maxDepthWithinDirectory := <-resultChannel
	fmt.Printf(
		"%s within %s directory has the highest depth of %d\n",
		maxDepthWithinDirectory.file,
		directory,
		maxDepthWithinDirectory.level,
	)
}

func forkIfNeeded(path string, info os.FileInfo, wg *sync.WaitGroup, results chan CodeDepth) {
	// if file is a directory or is not a .go file then don't process it
	if info.IsDir() || !strings.HasSuffix(path, ".go") {
		return
	}

	wg.Add(1)
	go func() {
		codeDepthForCurrentFile := deepestNestedBlock(path)
		results <- codeDepthForCurrentFile
		wg.Done()
	}()

}

func deepestNestedBlock(filename string) CodeDepth {
	code, _ := os.ReadFile(filename)
	maxDepth := 0
	currentDepth := 0

	for _, c := range code {
		if c == '{' {
			currentDepth += 1
			maxDepth = int(
				math.Max(
					float64(maxDepth),
					float64(currentDepth)),
			)
		} else if c == '}' {
			currentDepth -= 1
		}
	}
	return CodeDepth{
		file:  filename,
		level: maxDepth,
	}
}

func joinResults(codeDepths chan CodeDepth) chan CodeDepth {
	maxFileInDirectoryDepth := make(chan CodeDepth)

	go func() {
		max := CodeDepth{"", 0}
		for currentFileDepth := range codeDepths {
			if currentFileDepth.level > max.level {
				max = currentFileDepth
			}
		}
		/*
		 only sending the maxDepth entry when codeDepths channel is closed and
		 the above loop is terminated, signaling that we've gotten all depths for all file entries
		*/
		maxFileInDirectoryDepth <- max
	}()

	return maxFileInDirectoryDepth
}
