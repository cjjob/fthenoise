package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

type Message struct {
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}

type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

type Work struct {
	Title     string   `json:"title"`
	File      string   `json:"file"`
	Sentences []string `json:"sentences"`
}

type ReadPageData struct {
	WorksData template.JS
}

var works = []Work{
	{Title: "Pascal's Pensées", File: "texts/Pascal's Pensées.txt"},
	{Title: "Thus Spake Zarathustra", File: "texts/Thus Spake Zarathustra: A Book for All and None.txt"},
}

var parsedWorks = make(map[string]Work)

func init() {
	// Load and parse all works on startup
	for _, work := range works {
		sentences, err := loadAndParseWork(work.File)
		if err != nil {
			log.Printf("Warning: Failed to load %s: %v", work.File, err)
			continue
		}
		work.Sentences = sentences
		parsedWorks[work.File] = work
		log.Printf("Loaded %s: %d sentences", work.Title, len(sentences))
	}
}

func loadAndParseWork(filename string) ([]string, error) {
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
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/api/hello", helloHandler)
	http.HandleFunc("/breathe", breatheHandler)
	http.HandleFunc("/read", readHandler)

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
		log.Fatal("Server failed to start:", err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := `
<!DOCTYPE html>
<html>
<head>
	<title>Go Web Server</title>
	<style>
		body {
			font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
			max-width: 800px;
			margin: 50px auto;
			padding: 20px;
			background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
			color: white;
		}
		.container {
			background: rgba(255, 255, 255, 0.1);
			backdrop-filter: blur(10px);
			border-radius: 20px;
			padding: 40px;
			box-shadow: 0 8px 32px 0 rgba(31, 38, 135, 0.37);
		}
		h1 { margin-top: 0; }
		.endpoint {
			background: rgba(255, 255, 255, 0.2);
			padding: 15px;
			margin: 10px 0;
			border-radius: 10px;
			font-family: 'Courier New', monospace;
		}
		a {
			color: #fff;
			text-decoration: none;
			font-weight: bold;
		}
		a:hover { text-decoration: underline; }
	</style>
</head>
<body>
	<div class="container">
		<h1>🚀 Go Web Server</h1>
		<p>Welcome to your Go web server!</p>
		<h2>Available Endpoints:</h2>
		<div class="endpoint">
			<strong>GET</strong> <a href="/health">/health</a> - Health check endpoint
		</div>
		<div class="endpoint">
			<strong>GET</strong> <a href="/api/hello">/api/hello</a> - Hello API endpoint
		</div>
	</div>
</body>
</html>
	`
	fmt.Fprint(w, html)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
	}
	json.NewEncoder(w).Encode(response)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	message := Message{
		Message: "Hello from Go web server!",
		Time:    time.Now(),
	}
	json.NewEncoder(w).Encode(message)
}

func breatheHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/breathe" {
		http.NotFound(w, r)
		return
	}

	html := `In progress...`
	fmt.Fprint(w, html)
}

func readHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/read" {
		http.NotFound(w, r)
		return
	}

	// Convert parsed works to a map for JSON encoding
	worksMap := make(map[string]Work)
	for file, work := range parsedWorks {
		worksMap[file] = work
	}

	// Marshal to JSON and convert to template.JS for safe embedding
	worksJSON, err := json.Marshal(worksMap)
	if err != nil {
		http.Error(w, "Failed to prepare data", http.StatusInternalServerError)
		return
	}

	data := ReadPageData{
		WorksData: template.JS(worksJSON),
	}

	tmpl, err := template.ParseFiles("templates/read.html")
	if err != nil {
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}
