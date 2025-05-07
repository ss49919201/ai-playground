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

	// File server for the 'static' directory
	fs := http.FileServer(http.Dir("./static"))

	// Serve static files (CSS, JS, images etc.) from /static/ path
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Serve index.html at the root path "/"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// If the path is exactly "/", serve index.html from the static dir
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "./static/index.html")
			return
		}
		// Otherwise, try to serve as a short URL redirect
		app.handleRedirect(w, r)
	})

	// API endpoint for shortening
	http.HandleFunc("/shorten", app.handleShorten)
	// Note: The redirect handler is now part of the "/" handler logic above.

	port := ":8080"
	fmt.Printf("Server starting on http://localhost%s\n", port)
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
// This is now called by the main "/" handler if the path is not "/".
func (app *App) handleRedirect(w http.ResponseWriter, r *http.Request) {
	// Basic check to avoid trying to redirect favicon.ico etc. as short IDs
	// A more robust check might involve validating the ID format.
	if r.URL.Path == "/favicon.ico" {
		http.NotFound(w, r)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/")
	if id == "" { // Should not happen if called from "/" handler correctly
		http.NotFound(w, r)
		return
	}

	longURL, ok := app.store.Get(id)
	if !ok {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, longURL, http.StatusFound) // 302 Found redirect
}
