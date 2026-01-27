package api

import (
	"html/template"
	"net/http"
	"path/filepath"
)

func ExampleHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		filepath.Join("web", "templates", "base.html"),
		filepath.Join("web", "templates", "example.html"),
	)
	if err != nil {
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, "example", nil); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}
