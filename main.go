package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

// A struct to hold the front matter and content
type BlogPost struct {
	Title      string
	Excerpt    string
	Date       string
	Content    template.HTML
	Filename   string
	ParsedDate time.Time
}

func main() {
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// Serve about.html on the /about route
	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/about.html")
	})
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("GET /posts/{slug}", PostHandler)
	http.HandleFunc("/subscribe", SubscribeHandler)

	port := 8080
	address := fmt.Sprintf(":%d", port) // :8080

	log.Printf("Server is listening on port %d...\n", port)

	if err := http.ListenAndServe(address, nil); err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
}
