package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	// getting only the filenames from the command
	searchTerm := os.Args[1]
	directory := os.Args[2]

	go searchDir(searchTerm, directory)

	time.Sleep(3 * time.Second)
}

func searchDir(searchTerm string, directory string) {
	files, err := os.ReadDir(directory)
	if err != nil {
		log.Fatalf("Unable to read directory: %v, err: %v", directory, err)
	}

	for _, file := range files {
		fullPath := filepath.Join(directory, file.Name())
		// handle recursive directory search
		if file.IsDir() {

			go searchDir(searchTerm, fullPath)
		} else { // handle file search
			go searchFile(searchTerm, fullPath)
		}
	}

}

func searchFile(searchTerm string, filename string) {
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("Unable to read file: %v. err: %v\n", filename, err)
		return
	}

	stringifiedFileContent := string(fileContent)

	if strings.Contains(stringifiedFileContent, searchTerm) {
		log.Printf("file: %v contains a match for %v\n", filename, searchTerm)
	} else {
		log.Printf("file: %v does not contains a match for %v\n", filename, searchTerm)
	}
}

// go run grepdirrec.go {searchTerm} {directory}

// go run grepdirrec.go a ../../../common-files
