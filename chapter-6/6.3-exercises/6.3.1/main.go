package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}
	matchedFilesMutex := sync.Mutex{}
	matchedFiles := []string{}
	wg.Add(1)
	go fileSearch(os.Args[1], os.Args[2], &wg, &matchedFiles, &matchedFilesMutex)
	wg.Wait()

	matchedFilesMutex.Lock()
	sort.Strings(matchedFiles)
	fmt.Println(matchedFiles)
	matchedFilesMutex.Unlock()
	fmt.Printf("done searching for %v in directory %v", os.Args[1], os.Args[2])
}

func fileSearch(
	dir string,
	filename string,
	wg *sync.WaitGroup,
	results *[]string,
	mutex *sync.Mutex,
) {
	files, _ := os.ReadDir(dir)
	for _, file := range files {
		fpath := filepath.Join(dir, file.Name())
		if strings.Contains(file.Name(), filename) {
			mutex.Lock()
			*results = append(*results, fpath)
			mutex.Unlock()
		}
		if file.IsDir() {
			wg.Add(1)
			go fileSearch(fpath, filename, wg, results, mutex)
		}
	}
	wg.Done()
}
