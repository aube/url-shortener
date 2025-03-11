package handlers

import (
	"net/http"

	"github.com/aube/url-shortener/internal/logger"
)

type StorageGet interface {
	Get(key string) (value string, ok bool)
}

func HandlerID(store StorageGet) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		if id == "" {
			http.Error(w, "ID must be specified", http.StatusBadRequest)
			return
		}

		logger.Println("Requested ID:", id)

		url, ok := store.Get(id)
		if url == "" || !ok {
			http.Error(w, "URL not found", http.StatusBadRequest)
			return
		}

		logger.Println("URL:", url)

		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusTemporaryRedirect)

	}
}
