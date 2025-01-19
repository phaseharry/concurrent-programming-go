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

	files, err := os.ReadDir(directory)
	if err != nil {
		log.Fatalf("Unable to read directory: %v, err: %v", directory, err)
	}

	for _, file := range files {
		if !file.IsDir() {
			fullPath := filepath.Join(directory, file.Name())
			go searchFile(searchTerm, fullPath)
		}
	}

	time.Sleep(2 * time.Second)
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

// go run grepdir.go {searchTerm} {directory}

// go run grepdir.go a ../../../common-files
