package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func ExtractFrontMatter(markdownContent string) (map[string]string, string) {
	frontMatterPattern := regexp.MustCompile(`(?s)---(.*?)---(.*)`)
	matches := frontMatterPattern.FindStringSubmatch(markdownContent)

	frontMatter := make(map[string]string)
	var content string

	if len(matches) == 3 {
		for _, line := range strings.Split(matches[1], "\n") {
			if line == "" {
				continue
			}
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				// Remove surrounding quotes
				value = strings.Trim(value, `"`) // Trims just double quotes
				frontMatter[key] = value
			}
		}
		content = matches[2]
	} else {
		content = markdownContent // No front matter, entire content is Markdown
	}

	return frontMatter, content
}

func MdToHTML(md []byte) []byte {
	// Create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	var printAst = false
	// For debugging, not necessary in production
	if printAst {
		fmt.Print("--- AST tree:\n")
		ast.Print(os.Stdout, doc)
		fmt.Print("\n")
	}

	// Create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

func SubscribeEmailToMailchimp(email string) error {

	mailchimpDataCenter := os.Getenv("MAILCHIMP_DATA_CENTER")
	mailchimpApiKey := os.Getenv("MAILCHIMP_API_KEY")
	mailchimpListId := os.Getenv("MAILCHIMP_LIST_ID")

	if mailchimpDataCenter == "" || mailchimpApiKey == "" || mailchimpListId == "" {
		fmt.Println("Mailchimp environment variables are not set")
		return fmt.Errorf("mailchimp environment variables are not set")
	}

	apiUrl := fmt.Sprintf("https://%s.api.mailchimp.com/3.0/lists/%s/members", mailchimpDataCenter, mailchimpListId)

	payload := map[string]interface{}{
		"email_address": email,
		"status":        "subscribed", // Or "pending" if you want double opt-in
	}

	jsonPayload, err := json.Marshal(payload) // Convert to json
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+BasicAuth("", mailchimpApiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		// Handle error responses from Mailchimp
		return fmt.Errorf("mailchimp responded with status code: %d", resp.StatusCode)
	}

	return nil
}

// Helper function for HTTP Basic Auth
func BasicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
