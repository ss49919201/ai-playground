package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// App holds the application dependencies
type App struct {
	store URLStore
}

func main() {
	store := NewInMemoryURLStore()
	app := &App{store: store}

	http.HandleFunc("/shorten", app.handleShorten)
	http.HandleFunc("/", app.handleRedirect)

	port := ":8080"
	fmt.Printf("Server starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

// handleShorten handles requests to shorten a URL.
func (app *App) handleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	longURL := string(body)
	if longURL == "" {
		http.Error(w, "URL parameter is missing", http.StatusBadRequest)
		return
	}

	// Basic validation: Check if it looks like a URL
	_, err = url.ParseRequestURI(longURL)
	if err != nil {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	id := app.store.GenerateID()
	err = app.store.Set(id, longURL)
	if err != nil {
		http.Error(w, "Failed to store URL", http.StatusInternalServerError)
		return
	}

	// Construct the short URL (adjust scheme/host as needed)
	// Assuming server runs on localhost:8080 for now
	shortURL := fmt.Sprintf("http://localhost%s/%s", r.Host, id)
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, shortURL)
}

// handleRedirect handles requests to redirect short URLs to original URLs.
func (app *App) handleRedirect(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" || r.URL.Path == "/shorten" {
		// Handle root or the shorten path itself if needed, or return 404/specific page
		http.NotFound(w, r)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/")
	longURL, ok := app.store.Get(id)
	if !ok {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, longURL, http.StatusFound) // 302 Found redirect
}
