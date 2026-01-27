package typesgo

import (
	"html/template"
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

type Document struct {
	Title     string   `json:"title"`
	File      string   `json:"file"`
	Sentences []string `json:"sentences"`
}

type ReadPageData struct {
	DocumentsData template.JS
}
