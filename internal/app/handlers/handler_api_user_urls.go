package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	appErrors "github.com/aube/url-shortener/internal/app/apperrors"
	"github.com/aube/url-shortener/internal/logger"
)

type StorageList interface {
	List(c context.Context) (map[string]string, error)
}

func HandlerAPIUserUrls(ctx context.Context, store StorageList, baseURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.WithContext(ctx)
		w.Header().Set("Content-Type", "application/json")

		// userID := r.Context().Value(ctxkeys.UserIDKey)
		// ctx := context.WithValue(ctx, ctxkeys.UserIDKey, userID)

		memData, err := store.List(r.Context())

		var herr *appErrors.HTTPError
		if errors.As(err, &herr) {
			w.WriteHeader(herr.Code)
			return
		}

		if len(memData) == 0 {
			w.WriteHeader(204)
			return
		}

		json, err := getJSON(memData, baseURL)

		if err != nil {
			log.Error("getJSON", "err", err)
		}

		fmt.Fprintf(w, `%s`, string(json))

		w.WriteHeader(http.StatusOK)
	}
}

type JSONItem struct {
	Hash string `json:"short_url"`
	URL  string `json:"original_url"`
}

func getJSON(memData map[string]string, baseURL string) ([]byte, error) {

	var jsonData []JSONItem

	for k, v := range memData {
		item := JSONItem{Hash: baseURL + "/" + k, URL: v}
		jsonData = append(jsonData, item)
	}

	return json.Marshal(jsonData)
}
