package main

import (
	"fmt"
	"log"
	"os"
	"word-count/fileutil"
	"word-count/wordcounter"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <directory_path>")
		return
	}

	rootDir := os.Args[1]
	fmt.Printf("Root Directory: %v\n", rootDir)
	fmt.Println("\n----------Line count of txt files--------")

	textFiles, err := fileutil.ScanForTextFiles(rootDir)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("\n--------------Word Frequency-------------")

	wordCount := make(map[string]int)
	for _, file := range textFiles {
		wordcounter.CountWordsInFile(file, wordCount)
	}
	topN := 5
	wordcounter.PrintTopNWords(wordCount, topN)
}
