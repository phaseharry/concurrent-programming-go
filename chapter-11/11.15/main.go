package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

/*
demo of deadlock happening in go channels:

When we pass in initial directory, it will get sent to the dirsChannel and the handleDirectories goroutine will first process it.
That will send all of its sub directory/files to the files channel. Since the channels are not buffered, it can only send 1 path
at a time and then be blocked until the filesChannel is unblocked again / get its path processed in the handleFiles goroutine.
If a subdirectory was passed in to the handleFiles goroutine then it will lead to a deadlock as, the handlesFiles goroutine will
send that path back to handleDirectories goroutine through the dirsChannel. The handleDirectories goroutine is blocked because
the directory path sent to handlesFile through the filesChannel can't be processed and free up since it wasn't able to finish processing it,
while the handleFile goroutine is blocked because the dirs channel is still not finished with processing the initial directory path.

This circular dependency is causing a deadlock and neither goroutines can continue processing once a sub-directory is encountered. While
increasing the buffer for the channels might reduce the chance of this happening, if there is N subdirectories within a directory
and your buffers are N-1 then it will eventually lead to a deadlock
*/
func main() {
	filesChannel := make(chan string)
	dirsChannel := make(chan string)
	go handleFiles(filesChannel, dirsChannel)
	go handleDirectories(dirsChannel, filesChannel)
	dirsChannel <- os.Args[1]
	time.Sleep(60 * time.Second)
}

/*
consumes messages from a directory channel and for each file/path
within that directory, send the path to the files channel so recursively
go through all files/subdirectories within a directory
*/
func handleDirectories(dirs <-chan string, files chan<- string) {
	for fullpath := range dirs {
		fmt.Println("Reading all files from", fullpath)
		filesInDir, _ := os.ReadDir(fullpath)
		fmt.Printf("Pushing %d files from %s\n", len(filesInDir), fullpath)
		for _, file := range filesInDir {
			files <- filepath.Join(fullpath, file.Name())
		}
	}
}

/*
consumes files from files channel and check if it's a directory.
if it is a directory then push that path to the dirs channel so the
handleDirectories goroutine can process it.
if it's a file just print out its info
*/
func handleFiles(files chan string, dirs chan string) {
	for path := range files {
		file, _ := os.Open(path)
		fileInfo, _ := file.Stat()
		if fileInfo.IsDir() {
			fmt.Printf("Pushing %s directory\n", fileInfo.Name())
			dirs <- path
		} else {
			fmt.Printf("File %s, size %dMB, last modified: %s\n",
				fileInfo.Name(),
				fileInfo.Size()/(1024*1024),
				fileInfo.ModTime().Format("15:04:05"),
			)
		}
	}
}
