package wordcounter

import (
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"
)

type WordCount struct {
	Word  string
	Count int
}

type WordCountHeap []WordCount

func (h WordCountHeap) Len() int           { return len(h) }
func (h WordCountHeap) Less(i, j int) bool { return h[i].Count < h[j].Count }
func (h WordCountHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *WordCountHeap) Push(x interface{}) {
	*h = append(*h, x.(WordCount))
}

func (h *WordCountHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func CountWordsInFile(filename string, wordCount map[string]int) {
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Error opening file %s: %v\n", filename, err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		word := strings.ToLower(scanner.Text())
		word = strings.TrimFunc(word, func(r rune) bool { return !unicode.IsLetter(r) })
		if word != "" {
			wordCount[word]++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error scanning file %s: %v\n", filename, err)
	}
}

func PrintTopNWords(wordCount map[string]int, topN int) {

	priorityQueue := make(WordCountHeap, 0, topN)
	heap.Init(&priorityQueue)

	for word, count := range wordCount {
		if len(priorityQueue) < topN {
			heap.Push(&priorityQueue, WordCount{word, count})
		} else if count > priorityQueue[0].Count {
			heap.Pop(&priorityQueue)
			heap.Push(&priorityQueue, WordCount{word, count})
		}
	}

	// Extract and print the top N words
	fmt.Printf("Top %v most frequent words:\n", topN)
	for len(priorityQueue) > 0 {
		wc := heap.Pop(&priorityQueue).(WordCount)
		fmt.Printf("%s: %d\n", wc.Word, wc.Count)
	}
}
