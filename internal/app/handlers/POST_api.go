package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/aube/url-shortener/internal/app/hasher"
	"github.com/aube/url-shortener/internal/logger"
)

type StorageSet interface {
	Set(c context.Context, key string, value string) error
}

func HandlerAPI(store StorageSet, baseURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logger.WithContext(ctx)

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

		originalURL := readURLFromJSON(body)
		hash := hasher.CalcHash(originalURL)

		w.Header().Set("Content-Type", "application/json")
		httpStatus := http.StatusCreated

		err = store.Set(r.Context(), hash, string(originalURL))

		if err != nil {
			httpStatus = http.StatusConflict
		}
		w.WriteHeader(httpStatus)

		shortURL := baseURL + "/" + hash
		fmt.Fprintf(w, `{"result":"%s"}`, shortURL)

		log.Debug("URL:", shortURL, httpStatus)
	}
}
