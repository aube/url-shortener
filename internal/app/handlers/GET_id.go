package handlers

import (
	"context"
	"net/http"

	"github.com/aube/url-shortener/internal/logger"
)

type StorageGet interface {
	Get(c context.Context, key string) (value string, ok bool)
}

func HandlerID(store StorageGet) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logger.WithContext(ctx)

		id := r.PathValue("id")

		if id == "" {
			http.Error(w, "ID must be specified", http.StatusBadRequest)
			return
		}

		log.Debug("HandlerID", "id", id)

		url, ok := store.Get(ctx, id)
		if !ok {
			http.Error(w, "URL not found", http.StatusBadRequest)
			return
		}
		if url == "" && ok {
			http.Error(w, "URL deleted", http.StatusGone)
			return
		}

		log.Debug("HandlerID", "url", url)

		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusTemporaryRedirect)

	}
}
