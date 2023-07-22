package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
)

type Shortener struct {
	URLMap map[string]string
}

func (s *Shortener) GenerateShortURL(longURL string) string {
	hash := sha256.Sum256([]byte(longURL))
	shortURL := base64.URLEncoding.EncodeToString(hash[:8])
	s.URLMap[shortURL] = longURL
	return shortURL
}

func (s *Shortener) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	shortURL := r.URL.Path[1:]
	longURL, exists := s.URLMap[shortURL]
	if !exists {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, longURL, http.StatusMovedPermanently)
}

func main() {
	shortener := &Shortener{URLMap: make(map[string]string)}
	http.HandleFunc("/shorten", func(w http.ResponseWriter, r *http.Request) {
		longURL := r.URL.Query().Get("url")
		if longURL == "" {
			http.Error(w, "Missing Url Parameter", http.StatusBadRequest)
		}
		shortURL := shortener.GenerateShortURL(longURL)
		fmt.Fprintf(w, "Short URL: %s\n", shortURL)
	})

	http.HandleFunc("/", shortener.RedirectHandler)
	log.Println("Starting the server at 8080..")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
