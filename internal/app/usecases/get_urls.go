package usecases

import (
	"context"
	"encoding/json"

	"github.com/aube/url-shortener/internal/app/store"
	"github.com/aube/url-shortener/internal/logger"
)

// StorageList interface
type StorageList interface {
	List(c context.Context) (map[string]string, error)
}

func GetURLS(ctx context.Context, baseURL string) ([]byte, error) {
	store := store.NewStore()
	log := logger.WithContext(ctx)

	urls, err := store.List(ctx)
	if err != nil {
		log.Error("GetURLS", "err", err)
		return nil, err
	}

	json, err := getJSON(urls, baseURL)

	if err != nil {
		log.Error("GetURLS", "err2", err)
		return nil, err
	}

	return json, nil
}

// JSONItem struct
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
