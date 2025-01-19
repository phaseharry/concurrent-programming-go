package main

import (
	"log"
	"os"
	"time"
)

func main() {
	// getting only the filenames from the command
	filenames := os.Args[1:]

	for _, filename := range filenames {
		go readFile(filename)
	}

	time.Sleep(1 * time.Second)
}

func readFile(filename string) {
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("Unable to read file: %v. err: %v", filename, err)
	}

	stringifiedFileContent := string(fileContent)
	log.Println(stringifiedFileContent)

}

// go run catfiles.go ~/vim-notes.txt ~/test.sql
