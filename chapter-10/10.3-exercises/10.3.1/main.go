package main

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path"
	"sync"

	example10_1 "github.com/phaseharry/concurrent-programming-go/chapter-10/10.1"
)

/*
Example of Loop-carried dependence using waitGroups instead of the channel
implementation in example 10.5
*/

func main() {
	directory := os.Args[1]
	files, _ := os.ReadDir(directory)
	directorySha := sha256.New()

	var prevSignal, currentSignal *sync.WaitGroup

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		currentSignal = &sync.WaitGroup{}
		/*
		   Waiting for only 1 done call indicating that the current file's hash has been
		   appended to directory hash and the next file can do the same once it's ready.
		*/
		currentSignal.Add(1)

		go func(prev, current *sync.WaitGroup, filename string) {
			filePath := path.Join(directory, file.Name())
			fileHash := example10_1.FileHash(filePath)
			/*
				blocks current file from proceeding any further once it calculates its hash value.
				prevent it from appending to directory hash until all files before it in the directory
				has been added to the directory hash and signal the current file that it can add its hash.
				in this approach, we wait on the goroutine processing the previous file to add its fileHash
				to the directory hash before we add current file hash to directory hash. Once we add this file's hash
				to directory hash, we will call done on the current waitGroup, signaling the next file in the chain
				can append its hash to directory hash if it's ready
			*/
			if prev != nil {
				prev.Wait()
			}
			directorySha.Write(fileHash)
			current.Done()
		}(prevSignal, currentSignal, file.Name())

		// sent the prevSignal to currentSignal so the next file to be processed will be linked / chained properly
		prevSignal = currentSignal
	}

	// waits for the last file to be processed and have its hash appended to directory hash
	currentSignal.Wait()
	fmt.Printf("directory: %s, hash: %x\n", directory, directorySha.Sum(nil))
}
