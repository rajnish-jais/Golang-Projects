package fileutil

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func ScanForTextFiles(rootDir string) ([]string, error) {
	var textFiles []string
	filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && isTextFile(path) {
			lineCount, err := countLinesInFile(path)
			if err != nil {
				fmt.Printf("Error counting lines in %s: %v\n", path, err)
			} else {
				fmt.Printf("File: %s, Lines: %d\n", path, lineCount)
			}
			textFiles = append(textFiles, path)
		}
		return nil
	})
	return textFiles, nil
}

func isTextFile(filename string) bool {
	ext := filepath.Ext(filename)
	return strings.EqualFold(ext, ".txt")
}

func countLinesInFile(filename string) (int, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return 0, err
	}

	lines := strings.Split(string(data), "\n")
	return len(lines), nil
}
