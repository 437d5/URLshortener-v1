package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const alphabet = "abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ123456789"

const alphabetLen = len(alphabet)

const shortLenURL = 6

type URLShortener struct {
	urls map[string]string
}

func (shorten *URLShortener) HandleShortener(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "text/html")
		responseHTML := fmt.Sprintf(`<h2>URL Shortener</h2> <p>Input the URL you want to be short.</p>
		<form method="post" action="/shorten">
		<input type="text" name="url" placeholder="Enter a URL">
		<input type="submit" value="Shorten">
		</form>`)
		fmt.Fprintf(w, responseHTML)
	} else {

		fullURL := r.FormValue("url")
		if fullURL == "" {
			http.Error(w, "No URL parameter was passed", http.StatusBadRequest)
			return
		}

		shortToken := shortFullURL(alphabet, shortLenURL, alphabetLen)
		shorten.urls[shortToken] = fullURL

		shortURL := fmt.Sprintf("http://localhost:8080/short/%s", shortToken)

		w.Header().Set("Content-Type", "text/html")
		responseHTML := fmt.Sprintf(`
			<h2>URL Shortener</h2>
			<p>Original URL: %s</p>
			<p>Shortened URL: <a href="%s">%s</a></p>
			<form method="post" action="/shorten">
				<input type="text" name="url" placeholder="Enter a URL">
				<input type="submit" value="Shorten">
			</form>
		`, fullURL, shortURL, shortURL)
		fmt.Fprintf(w, responseHTML)
	}

}

func (shorten *URLShortener) HandlerRediret(w http.ResponseWriter, r *http.Request) {
	shortToken := r.URL.Path[len("/short/"):]
	if shortToken == "" {
		http.Error(w, "Shortened URL is missing", http.StatusBadRequest)
		return
	}

	fullURL, found := shorten.urls[shortToken]
	if !found {
		http.Error(w, "Shortened URL not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, fullURL, http.StatusSeeOther)
}

func shortFullURL(alphabet string, shortLenURL int, alphabetLen int) string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	shortToken := make([]byte, shortLenURL)
	for i := range shortToken {
		shortToken[i] = alphabet[rand.Intn(alphabetLen)]
	}
	return string(shortToken)
}

func main() {
	shorten := &URLShortener{
		urls: make(map[string]string),
	}

	http.HandleFunc("/shorten", shorten.HandleShortener)
	http.HandleFunc("/short/", shorten.HandlerRediret)

	fmt.Println("URL Shortener is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
