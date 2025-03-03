package handlers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/aube/url-shortener/internal/app/hasher"
	"github.com/aube/url-shortener/internal/app/store"
	"github.com/aube/url-shortener/internal/logger"
)

func HandlerAPI(baseURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		MemoryStore := store.NewMemoryStore()

		if r.Body == nil || r.ContentLength == 0 {
			http.Error(w, "Request body is empty", http.StatusBadRequest)
			return
		}

		// Read the entire body content
		body, err := io.ReadAll(r.Body)

		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		originalURL := readURLFromJSON(body)
		hash := hasher.CalcHash(originalURL)

		MemoryStore.Set(hash, string(originalURL))
		MemoryStore.Get(hash)

		shortURL := baseURL + "/" + hash

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprintf(w, `{"result":"%s"}`, shortURL)

		logger.Println("URL:", shortURL, http.StatusCreated)
	}
}
