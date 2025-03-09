package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	// taking arguments from the cli when running "ex. go run main.go directory filename"
	go fileSearch(os.Args[1], os.Args[2], &wg)
	wg.Wait()
	fmt.Printf("done searching for %v in directory %v", os.Args[1], os.Args[2])
}

/*
function used to search for a file within a directory.
supports nested directory searching, it will search through
every subdirectory to see if the file we're searching for is there.
recursively spawns new go routines when a directory is detected
and that go routine will be responsible for checking if the
file exists within it and spawn goroutines to do the same for that
directory's subdirectories.
*/
func fileSearch(dir string, filename string, wg *sync.WaitGroup) {
	// getting all files within current directory
	files, _ := os.ReadDir(dir)

	for _, file := range files {
		fpath := filepath.Join(dir, file.Name())
		/*
		   checking if the current file is the filename we're looking for.
		   if so print it, but continue searching
		*/
		if strings.Contains(file.Name(), filename) {
			fmt.Println(fpath)
		}
		/*
		   if the current file is a directory, spawn a new goroutine to check its contents.
		   increment the WaitGroup's counter by 1 so the main goroutine will wait until
		   all the files and subdirectories within the current directory has been searched.
		*/
		if file.IsDir() {
			wg.Add(1)
			go fileSearch(fpath, filename, wg)
		}
	}
	// after we've searched through every entity for current directory, call Done()
	wg.Done()
}
