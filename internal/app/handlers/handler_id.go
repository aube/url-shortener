package handlers

import (
	"fmt"
	"net/http"

	"github.com/aube/url-shortener/internal/app/store"
)

func HandlerID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if id == "" {
		http.Error(w, "ID must be specified", http.StatusBadRequest)
		return
	}

	url := store.GetURLHash(id)

	if url == "" {
		http.Error(w, "URL not found", http.StatusBadRequest)
		return
	}

	fmt.Println("ID:", id)
	fmt.Println("URL:", url)

	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)

}
