package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	example10_1 "github.com/phaseharry/concurrent-programming-go/chapter-10/10.1"
)

/*
Using loop-level parallelism to compute hashcodes of file.
Iterate through list of files and for each file, spin off a new
goroutine / thread to concurrently calculate hash code for each file.

this can only be done because the file hashing of a file is independent of the
file hashing another file so it can all be concurrently done
*/
func main() {
	directory := os.Args[1]
	files, _ := os.ReadDir(directory)
	waitGroup := sync.WaitGroup{}

	for _, file := range files {

		/*
			only handling files within current directory. will not recursively
			iterate through subdirectories and handling their files
		*/
		if file.IsDir() {
			continue
		}

		waitGroup.Add(1)

		// spinning a new goroutine for each file
		go func(filename string) {
			filePath := filepath.Join(directory, filename)
			fileHash := example10_1.FileHash(filePath)
			fmt.Printf("%s - %x\n", filename, fileHash)
			waitGroup.Done()
		}(file.Name())
	}

	// wait until all files have been hand;ed
	waitGroup.Wait()
}
