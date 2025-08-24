package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func main() {
	filesChannel := make(chan string)
	dirsChannel := make(chan string)
	go handleFiles(filesChannel, dirsChannel)
	go handleDirectories(dirsChannel, filesChannel)
	dirsChannel <- os.Args[1]
	time.Sleep(60 * time.Second)
}

/*
takes files/subdirs from the files channel (consuming from) and determine if the
entry is a file or subdirectory. If it's a file, print file info and stop nesting further.
If it's a subdirectory, publish the path to the dirs channel and let handleDirectories goroutine
process and publish back sub directories back to the handleFiles goroutine
*/
func handleFiles(files <-chan string, dirs chan<- string) {
	dirsToPush := make([]string, 0)
	appendAllDirs := func(path string) {
		file, _ := os.Open(path)
		fileInfo, _ := file.Stat()
		if fileInfo.IsDir() {
			fmt.Printf("Pushing %s directory\n", fileInfo.Name())
			dirsToPush = append(dirsToPush, path)
		} else {
			fmt.Printf("File %s, size %dMB, last modified: %s\n",
				fileInfo.Name(),
				fileInfo.Size()/(1024*1024),
				fileInfo.ModTime().Format("15:04:05"),
			)
		}
	}

	for {
		if len(dirsToPush) == 0 {
			appendAllDirs(<-files)
		} else {
			select {
			case fullpath := <-files:
				appendAllDirs(fullpath)
			case dirs <- dirsToPush[0]:
				dirsToPush = dirsToPush[1:]
			}
		}
	}
}

/*
takes directory / file paths from dirs channel (consuming from) and gets all files within that path and
send it to the files channel to be consumed and processed by the handlesFiles goroutine
*/
func handleDirectories(dirs <-chan string, files chan<- string) {
	toPush := make([]string, 0)
	appendAllFiles := func(path string) {
		fmt.Println("Reading all files from", path)
		filesInDir, _ := os.ReadDir(path)
		fmt.Printf("Pushung %d files from %s\n", len(filesInDir), path)
		for _, f := range filesInDir {
			toPush = append(toPush, filepath.Join(path, f.Name()))
		}
	}

	for {
		if len(toPush) == 0 {
			appendAllFiles(<-dirs)
		} else {
			select {
			case fullpath := <-dirs:
				appendAllFiles(fullpath)
			case files <- toPush[0]:
				toPush = toPush[1:]
			}
		}
	}
}
