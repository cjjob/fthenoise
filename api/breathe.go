package api

import (
	"fmt"
	"net/http"
)

func BreatheHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/breathe" {
		http.NotFound(w, r)
		return
	}

	html := `In progress...`
	fmt.Fprint(w, html)
}
