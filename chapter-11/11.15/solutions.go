package main

import (
	"fmt"
	"os"
	"path/filepath"
)

/*
this as a replacement fgor handleDirectories would prevent the dead lock
as it spins off a new goroutine to send it to the filesChannel without blocking
the handleDirectories goroutine. This breaks the circular wait between handleDirectories
and handleFiles so they won't block either other
*/
func handleDirectoriesSolution1(dirs <-chan string, files chan<- string) {
	for fullpath := range dirs {
		fmt.Println("Reading all files from", fullpath)
		filesInDir, _ := os.ReadDir(fullpath)
		fmt.Printf("Pushing %d files from %s\n", len(filesInDir), fullpath)
		for _, file := range filesInDir {
			go func(fp string) {
				files <- fp
			}(filepath.Join(fullpath, file.Name()))
		}
	}
}

/*
using a select statement to break the circular wait.

if there is nothing to push (toPush slice is empty) then we will block
until we get a directory path that can be consumed from the dirsChannel.

once there's a path that's coming from the dirsChannel, it will consume that path
and get all subdirectories/files within that path and add it to the toPush slice.
in the next iteration of the for loop, the else -> select case will take over
and either consume paths within the toPush until it's empty or append more subdirectories
from the dirsChannel into the toPush slice.
*/
func handleDirectoriesSolution2(dirs <-chan string, files chan<- string) {
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
