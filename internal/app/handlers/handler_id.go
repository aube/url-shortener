package handlers

import (
	"context"
	"net/http"

	"github.com/aube/url-shortener/internal/logger"
)

type StorageGet interface {
	Get(c context.Context, key string) (value string, ok bool)
}

func HandlerID(ctx context.Context, store StorageGet) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		if id == "" {
			http.Error(w, "ID must be specified", http.StatusBadRequest)
			return
		}

		logger.Println("Requested ID:", id)

		url, ok := store.Get(ctx, id)
		if url == "" || !ok {
			http.Error(w, "URL not found", http.StatusBadRequest)
			return
		}

		logger.Println("URL:", url)

		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusTemporaryRedirect)

	}
}
