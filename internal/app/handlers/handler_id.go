package handlers

import (
	"fmt"
	"net/http"

	"github.com/aube/url-shortener/internal/app/hashes"
)

func HandlerId(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		id := r.PathValue("id")

		if id == "" {
			http.Error(w, "ID must be specified", http.StatusBadRequest)
			return
		}

		url := hashes.GetURLHash(id)

		if url == "" {
			http.Error(w, "URL not found", http.StatusBadRequest)
			return
		}

		fmt.Println("ID:", id)
		fmt.Println("URL:", url)

		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusTemporaryRedirect)

	default:
		fmt.Println("Not served method:", r.Method)
	}
}
