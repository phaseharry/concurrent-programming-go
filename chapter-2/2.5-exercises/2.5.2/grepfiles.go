package main

import (
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	// getting only the filenames from the command
	searchTerm := os.Args[1]
	filenames := os.Args[2:]

	for _, filename := range filenames {
		go searchFile(searchTerm, filename)
	}

	time.Sleep(1 * time.Second)
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

// go run grepfiles.go a ~/vim-notes.txt ~/test.sql
