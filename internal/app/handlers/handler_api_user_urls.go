package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aube/url-shortener/internal/logger"
)

type StorageList interface {
	List() map[string]string
}

func HandlerAPIUserUrls(store StorageList, baseURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		memData := store.List()
		json := getJSON(memData, baseURL)
		fmt.Fprintf(w, `%s`, json)

		w.WriteHeader(http.StatusOK)
	}
}

type JSONItem struct {
	Hash string `json:"short_url"`
	URL  string `json:"original_url"`
}

func getJSON(memData map[string]string, baseURL string) string {
	var jsonData []JSONItem

	for k, v := range memData {
		item := JSONItem{Hash: baseURL + "/" + k, URL: v}
		jsonData = append(jsonData, item)
	}
	jsonBytes, err := json.Marshal(jsonData)

	if err != nil {
		logger.Infoln(err)
	}

	return string(jsonBytes)
}
