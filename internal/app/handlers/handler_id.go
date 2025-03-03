package handlers

import (
	"net/http"

	"github.com/aube/url-shortener/internal/app/store"
	"github.com/aube/url-shortener/internal/logger"
)

func HandlerID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		MemoryStore := store.NewMemoryStore()

		id := r.PathValue("id")

		if id == "" {
			http.Error(w, "ID must be specified", http.StatusBadRequest)
			return
		}

		logger.Println("Requested ID:", id)

		url, ok := MemoryStore.Get(id)
		if url == "" || !ok {
			http.Error(w, "URL not found", http.StatusBadRequest)
			return
		}

		logger.Println("URL:", url)

		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusTemporaryRedirect)

	}
}
