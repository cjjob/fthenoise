package main

import (
	"encoding/json"
	"fmt"
	"fthenoise/api"
	typesgo "fthenoise/types.go"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var ParsedDocuments = make(map[string]typesgo.Document)

func init() {
	// Load and parse all documents on startup
	for _, document := range documents {
		sentences, err := loadAndParseDocument(document.File)
		if err != nil {
			slog.Warn("Failed to load document",
				"file", document.File,
				"error", err,
			)
			continue
		}
		document.Sentences = sentences
		ParsedDocuments[document.File] = document
		slog.Info("Loaded document",
			"title", document.Title,
			"sentence_count", len(sentences),
		)
	}
	// Set the parsed documents in the api package
	api.ParsedDocuments = ParsedDocuments
}

func loadAndParseDocument(filename string) ([]string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	text := string(content)

	// Find the actual content (skip Project Gutenberg header)
	startMarker := "*** START OF THE PROJECT GUTENBERG EBOOK"
	startIdx := strings.Index(text, startMarker)
	if startIdx != -1 {
		// Find the end of the marker line and skip a few lines
		nextNewline := strings.Index(text[startIdx:], "\n")
		if nextNewline != -1 {
			text = text[startIdx+nextNewline+1:]
			// Skip a few more lines to get past headers
			for i := 0; i < 10; i++ {
				if idx := strings.Index(text, "\n"); idx != -1 {
					text = text[idx+1:]
				}
			}
		}
	}

	// Split into sentences using regex
	// Match sentence endings: . ! ? followed by space or newline
	sentenceRegex := regexp.MustCompile(`([.!?]+)\s+`)
	sentences := sentenceRegex.Split(text, -1)

	// Clean up sentences
	var cleanedSentences []string
	for _, sentence := range sentences {
		sentence = strings.TrimSpace(sentence)
		// Filter out very short sentences (likely artifacts) and empty strings
		if len(sentence) > 10 && !strings.HasPrefix(strings.ToLower(sentence), "produced by") {
			cleanedSentences = append(cleanedSentences, sentence)
		}
	}

	return cleanedSentences, nil
}

func main() {
	// Define routes
	http.HandleFunc("/", homeHandler)

	http.HandleFunc("/example", api.ExampleHandler)

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/breathe", api.BreatheHandler)
	http.HandleFunc("/read", api.ReadHandler)

	// Start server - read PORT from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}
	fmt.Printf("🚀 Server starting on http://localhost%s\n", port)
	fmt.Println("Available endpoints:")
	fmt.Println("  GET  /           - Home page")
	fmt.Println("  GET  /health     - Health check")
	fmt.Println("  GET  /api/hello  - Hello API endpoint")
	fmt.Println("  GET  /breathe    - Breathe endpoint")
	fmt.Println("  GET  /read       - Read endpoint")

	if err := http.ListenAndServe(port, nil); err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles(
		filepath.Join("web", "templates", "base.html"),
		filepath.Join("web", "templates", "home.html"),
	)
	if err != nil {
		slog.Error("Failed to load templates", "error", err)
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, "home", nil); err != nil {
		slog.Error("Failed to render template", "error", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := typesgo.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
	}
	json.NewEncoder(w).Encode(response)
}
