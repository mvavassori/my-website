package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Read all files in the 'posts' directory
	files, err := os.ReadDir("posts")
	if err != nil {
		log.Printf("Error reading posts directory: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Initialize an empty slice to hold the blog posts
	var posts []BlogPost
	// Loop through the files and extract front matter and content
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".md") {
			// Read the Markdown file
			markdownBytes, err := os.ReadFile(filepath.Join("posts", file.Name()))
			if err != nil {
				log.Println(err)
				continue // Skip this file
			}

			// Extract front matter
			frontMatter, _ := ExtractFrontMatter(string(markdownBytes))

			// Parse the date from the front matter
			parsedDate, err := time.Parse("2006-01-02", frontMatter["date"])
			if err != nil {
				log.Printf("Error parsing date: %v", err)
				continue
			}

			// Create a new BlogPost struct and append it to the slice
			posts = append(posts, BlogPost{
				Title:      frontMatter["title"],
				Excerpt:    frontMatter["excerpt"],
				Date:       frontMatter["date"],
				Filename:   strings.TrimSuffix(file.Name(), ".md"),
				ParsedDate: parsedDate,
			})
		}
	}

	// Sort the posts by date, newest first
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].ParsedDate.After(posts[j].ParsedDate)
	})

	// Parse and execute the template
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, posts)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	// Extract path parameter using new go's syntax
	postSlug := r.PathValue("slug")

	// Append '.md' to get the actual filename
	filename := postSlug + ".md"
	// Read the Markdown file
	markdownBytes, err := os.ReadFile(filepath.Join("posts", filename))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Extract front matter
	frontMatter, content := ExtractFrontMatter(string(markdownBytes))

	// Convert Markdown content to HTML
	htmlContent := MdToHTML([]byte(content))

	// Parse and execute the template with front matter and content
	tmpl, err := template.ParseFiles("templates/post.html")
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := BlogPost{
		Title:   frontMatter["title"],
		Excerpt: frontMatter["excerpt"],
		Date:    frontMatter["date"],
		Content: template.HTML(htmlContent),
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func SubscribeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestData struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Error decoding request", http.StatusBadRequest)
		return
	}

	fmt.Printf("Email submitted: %s\n", requestData.Email)

	err := SubscribeEmailToMailchimp(requestData.Email)
	if err != nil {
		// Handle error, possibly sending a different HTTP status or message to the client
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// After processing the subscription
	fmt.Fprintf(w, "Subscription successful for email: %s", requestData.Email)
}
