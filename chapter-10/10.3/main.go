package main

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path"

	example10_1 "github.com/phaseharry/concurrent-programming-go/chapter-10/10.1"
)

/*
Example of Loop-carried dependence. In this example, we want to calculate a directory's hash
to determine if any of the files within the directory has been changed. To do this, we calculate
the individual hashes of each file within the directory and write its bytes to the directory hash.
The order matters and the files would always been written in the same order as its previous if nothing
has changed. Go's ReadDir ensures this as the order will always be same (Directory Order).

This is easy in a sequential version as it is as simple as looping through the files and calculating the hash of the file
and just appending it. To do it concurrently requires more modification
*/

func main() {
	directory := os.Args[1]
	files, _ := os.ReadDir(directory)
	directorySha := sha256.New()
	/*
		since the ordering of the adding of the file hash to directory hash matters,
		we cannot just have a goroutine to operate each file hash and that's it. We want to
		maintain the ordering of adding the file hash to the directory hash.
		The adding of the hash result is dependent on the order of the files within the directory,
		but not the actual file hash computation. For each file, we can have a previous channel (prevSignal)
		that signals current goroutine when the previous goroutine has finished adding the file hash
		result to the directory hash, so that it can add its file hash to the directory hash.
		It also creates a new channel for itself and it will be used as a signal for the next goroutine (file)
		of when the current file hash has been added to the directory hash so it can add its file hash to directory hash
	*/
	var prevSignal, currentSignal chan int
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		/*
			currentSignal will be used by the current file to signal the next
			file in the directory that it can append its hash to the directory hash
			because all previous files have already been added.

			In the next file in the chain, this would be the `previousSignal`.
			Once it's unblocked by consuming the signal, it will add its hash to the directory hash
		*/
		currentSignal = make(chan int)
		go func(prev chan int, current chan int, filename string) {
			filePath := path.Join(directory, file.Name())
			fileHash := example10_1.FileHash(filePath)
			/*
				blocks current file from proceeding any further once it calculates its hash value.
				prevent it from appending to directory hash until all files before it in the directory
				has been added to the directory hash and signal the current file that it can add its hash.
			*/
			if prev != nil {
				<-prev
			}
			directorySha.Write(fileHash)
			current <- 0
		}(prevSignal, currentSignal, file.Name())

		// sent the prevSignal to currentSignal so the next file to be processed will be linked / chained properly
		prevSignal = currentSignal

	}

	// waits for the last file to be processed and have its hash appended to directory hash
	<-currentSignal
	fmt.Printf("directory: %s, hash: %x\n", directory, directorySha.Sum(nil))
}

// sequential version
func sequentialVersion() {
	dir := os.Args[1]
	files, _ := os.ReadDir(dir)

	// directory hash that all file hashes will be written to
	directorySha := sha256.New()

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := path.Join(dir, file.Name())
		fileHash := example10_1.FileHash(filePath)
		directorySha.Write(fileHash)
	}

	fmt.Printf("%s - %x\n", dir, directorySha.Sum(nil))
}
