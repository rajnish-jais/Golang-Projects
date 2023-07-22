package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Response struct {
	Url  string `json:"url"`
	Data string `json:"data"`
}

type Result struct {
	Result []Response `json:"result"`
	Error  error      `json:"error"`
}

func main() {
	var res Result
	urls, err := getUrlsFromFile()
	resp := make([]Response, len(urls))

	// crawl on each url
	for i, url := range urls {
		data, err := crawl(url)
		if err != nil {
			resp[i].Url = url
			resp[i].Data = err.Error()
		} else {
			resp[i].Url = url
			resp[i].Data = data
		}
	}

	// save the result and error
	res.Result = resp
	res.Error = err

	fmt.Println(prettyPrint(res))
}

func prettyPrint(i interface{}) string {
	s, err := json.MarshalIndent(i, "", "\t")
	if err != nil {
		fmt.Errorf("got error while marshalling the data: %v", err)
	}
	return string(s)
}

// crawl, it will fetch the url page content.
func crawl(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Errorf(" got error while fetching the url content: %v", err)
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	return string(body), nil
}

// getUrlsFromFile, reads from the url file and return a slice of strings
func getUrlsFromFile() ([]string, error) {
	urls := make([]string, 0)
	readFile, err := os.Open("url.txt")
	defer readFile.Close()
	if err != nil {
		fmt.Errorf("got error while reading from the file:%v", err)
		return urls, err
	}
	fileScanner := bufio.NewScanner(readFile)

	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		urls = append(urls, fileScanner.Text())
	}

	return urls, nil
}
