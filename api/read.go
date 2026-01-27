package api

import (
	"encoding/json"
	typesgo "fthenoise/types.go"
	"html/template"
	"net/http"
	"path/filepath"
)

// ParsedDocuments is set by main package after loading documents
var ParsedDocuments map[string]typesgo.Document

func ReadHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/read" {
		http.NotFound(w, r)
		return
	}

	// Convert parsed works to a map for JSON encoding
	documentsMap := make(map[string]typesgo.Document)
	for file, document := range ParsedDocuments {
		documentsMap[file] = document
	}

	// Marshal to JSON and convert to template.JS for safe embedding
	documentsJSON, err := json.Marshal(documentsMap)
	if err != nil {
		http.Error(w, "Failed to prepare data", http.StatusInternalServerError)
		return
	}

	data := typesgo.ReadPageData{
		DocumentsData: template.JS(documentsJSON),
	}

	tmpl, err := template.ParseFiles(
		filepath.Join("web", "templates", "base.html"),
		filepath.Join("web", "templates", "read.html"),
	)
	if err != nil {
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, "read", data); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}
