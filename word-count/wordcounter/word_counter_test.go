package wordcounter

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestCountWordsInFile(t *testing.T) {
	// Create a temporary test file with the sample text.
	fileContent := "This is a test text. This is another test."
	tempFile, err := createTempFileWithContent(fileContent)
	if err != nil {
		t.Fatalf("Error creating temporary test file: %v", err)
	}
	defer os.Remove(tempFile)

	wordCount := make(map[string]int)
	CountWordsInFile(tempFile, wordCount)

	expectedWordCount := map[string]int{
		"this":    2,
		"is":      2,
		"a":       1,
		"test":    2,
		"text":    1,
		"another": 1,
	}

	for word, count := range expectedWordCount {
		if wordCount[word] != count {
			t.Errorf("Expected count for word '%s' to be %d, but got %d", word, count, wordCount[word])
		}
	}
}

func createTempFileWithContent(content string) (string, error) {
	tempFile, err := os.CreateTemp("", "test-*.txt")
	if err != nil {
		return "", err
	}
	_, err = tempFile.WriteString(content)
	if err != nil {
		return "", err
	}
	return tempFile.Name(), nil
}

func TestPrintTopNWords(t *testing.T) {
	wordCount := map[string]int{
		"apple":  5,
		"banana": 3,
		"cherry": 8,
		"date":   2,
		"elder":  6,
		"fig":    7,
		"grape":  1,
		"honey":  4,
		"kiwi":   9,
		"lemon":  10,
		"mango":  11,
	}

	topN := 5
	expectedTopWords := []string{"mango: 11", "lemon: 10", "kiwi: 9", "fig: 7", "cherry: 8"}

	result := captureOutput(func() {
		PrintTopNWords(wordCount, topN)
	})

	for _, expectedWord := range expectedTopWords {
		if !strings.Contains(result, expectedWord) {
			t.Errorf("Expected output to contain '%s', but it did not.", expectedWord)
		}
	}
}

func captureOutput(f func()) string {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = oldStdout

	return string(out)
}
