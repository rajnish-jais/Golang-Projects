package fileutil

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// createTestFiles creates sample text files for testing purposes.
func createTestFiles(dir string) error {
	_ = os.MkdirAll(dir, os.ModePerm)
	fileContents := []string{"Hello, World!", "This is a test.", "Another line."}
	for i, content := range fileContents {
		filePath := filepath.Join(dir, fmt.Sprintf("test%d.txt", i))
		if err := ioutil.WriteFile(filePath, []byte(content), 0644); err != nil {
			return err
		}
	}
	return nil
}

func TestScanForTextFiles(t *testing.T) {
	tempDir := t.TempDir()
	if err := createTestFiles(tempDir); err != nil {
		t.Fatal(err)
	}

	files, err := ScanForTextFiles(tempDir)
	if err != nil {
		t.Fatalf("Error scanning for text files: %v", err)
	}

	expectedFiles := []string{
		filepath.Join(tempDir, "test0.txt"),
		filepath.Join(tempDir, "test1.txt"),
		filepath.Join(tempDir, "test2.txt"),
	}

	if len(files) != len(expectedFiles) {
		t.Errorf("Expected %d files, but got %d", len(expectedFiles), len(files))
	}

	for i := range expectedFiles {
		if files[i] != expectedFiles[i] {
			t.Errorf("Expected file %s, but got %s", expectedFiles[i], files[i])
		}
	}
}

func TestIsTextFile(t *testing.T) {
	testCases := []struct {
		filename   string
		isTextFile bool
	}{
		{"file.txt", true},
		{"file.TXT", true},
		{"file.doc", false},
		{"file.pdf", false},
	}

	for _, tc := range testCases {
		t.Run(tc.filename, func(t *testing.T) {
			if result := isTextFile(tc.filename); result != tc.isTextFile {
				t.Errorf("Expected %v for %s, but got %v", tc.isTextFile, tc.filename, result)
			}
		})
	}
}

func TestCountLinesInFile(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	fileContents := "Line 1\nLine 2\nLine 3"
	if err := ioutil.WriteFile(testFile, []byte(fileContents), 0644); err != nil {
		t.Fatal(err)
	}

	lineCount, err := countLinesInFile(testFile)
	if err != nil {
		t.Fatalf("Error counting lines in file: %v", err)
	}

	expectedLineCount := 3
	if lineCount != expectedLineCount {
		t.Errorf("Expected %d lines, but got %d", expectedLineCount, lineCount)
	}
}
