package usecases

import (
	"context"
	"encoding/json"

	"github.com/aube/url-shortener/internal/app/store"
	"github.com/aube/url-shortener/internal/logger"
)

// StorageList defines the interface for listing all URLs in storage.
type StorageList interface {
	// List returns a map of all hashes to their original URLs.
	// Returns an error if the operation fails.
	List(c context.Context) (map[string]string, error)
}

// GetURLS retrieves all stored URLs and returns them as JSON.
// The returned URLs are formatted with the provided baseURL.
// Returns the URLs as JSON bytes or an error if the operation fails.
func GetURLS(ctx context.Context, baseURL string) ([]byte, error) {
	store := store.NewStore()
	log := logger.WithContext(ctx)

	urls, err := store.List(ctx)
	if err != nil {
		log.Error("GetURLS", "err", err)
		return nil, err
	}

	if len(urls) == 0 {
		log.Info("GetURLS", "message", "no URLs found")
		return []byte(""), nil
	}

	json, err := getJSON(urls, baseURL)

	if err != nil {
		log.Error("GetURLS", "err2", err)
		return nil, err
	}

	return json, nil
}

// JSONItem represents a single URL entry in the JSON response.
type JSONItem struct {
	Hash string `json:"short_url"`    // The shortened URL with base
	URL  string `json:"original_url"` // The original URL
}

func getJSON(memData map[string]string, baseURL string) ([]byte, error) {

	var jsonData []JSONItem

	for k, v := range memData {
		item := JSONItem{Hash: baseURL + "/" + k, URL: v}
		jsonData = append(jsonData, item)
	}

	return json.Marshal(jsonData)
}
