package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type URL struct {
	ID           string    `json:"id"`
	OriginalURL  string    `json:"original_url"`
	ShortURL     string    `json:"short_url"`
	CreationDate time.Time `json:"creation_date"`
}

/*
d975485 --> {
ID: "d975485"
originalURL: "https://github.com/Prince-1501/"
shortURL: "d975485"
CreationDate: time.Now()}
*/

var urlDB = make(map[string]URL)

func generateShortURL(OriginalURL string) string {
	hasher := md5.New()
	hasher.Write([]byte(OriginalURL)) //it coverts original url string to a byte slice
	fmt.Println("hasher: ", hasher)
	data := hasher.Sum(nil)
	fmt.Println("hasher data is: ", data)
	hash := hex.EncodeToString(data)
	fmt.Println("encode to string:", hash[:8])

	return hash[:8]
}

func createURL(originalURL string) string {
	shortURL := generateShortURL(originalURL)
	id := shortURL // use the short url as the ID for simpllicity
	urlDB[id] = URL{
		ID:           id,
		OriginalURL:  originalURL,
		ShortURL:     shortURL,
		CreationDate: time.Now(),
	}
	return shortURL
}

func getURL(id string) (URL, error) {
	url, ok := urlDB[id]
	if !ok {
		return URL{}, errors.New(" url not found")
	}
	return url, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello, world!") // to write directly on the console of page
}

func ShortURLHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		URL string `json:"url"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "invallid request body", http.StatusBadRequest)
		return
	}
	shortURL_ := createURL(data.URL)
	//fmt.Fprintf(w, shortURL)
	response := struct {
		ShortURL string `json:"short_url"`
	}{ShortURL: shortURL_}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func redirectURLHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/redirect/"):]
	url, err := getURL(id)
	if err != nil {
		http.Error(w, "invalid request", http.StatusNotFound)
	}
	http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}

func main() {
	// fmt.Println("starting url shortener")
	// OriginalURL:= "https://github.com/Prince-1501/"
	// generateShortURL(OriginalURL)

	// register the handler function to handle all requests to the root URL ("/")
	http.HandleFunc("/", handler)
	http.HandleFunc("/shorten", ShortURLHandler)
	http.HandleFunc("/redirect/", redirectURLHandler)

	// server in golang
	fmt.Println("starting surver on port 3000...")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Println("error on starting server: ", err)
	}

}
