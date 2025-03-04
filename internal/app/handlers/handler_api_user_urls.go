package handlers

import (
	"fmt"
	"net/http"

	"github.com/aube/url-shortener/internal/app/store"
)

func HandlerAPIUserUrls(baseURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		MemoryStore := store.NewMemoryStore()

		w.Header().Set("Content-Type", "application/json")
		json := string(MemoryStore.JSON(baseURL))
		fmt.Fprintf(w, `%s`, json)

		w.WriteHeader(http.StatusOK)
	}
}
