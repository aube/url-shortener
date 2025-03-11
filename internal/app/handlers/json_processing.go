package handlers

import (
	"encoding/json"

	"github.com/aube/url-shortener/internal/logger"
)

type URLJson struct {
	URL string `json:"URL"`
}

func readURLFromJSON(body []byte) []byte {
	var jsonBody URLJson

	err := json.Unmarshal(body, &jsonBody)
	if err != nil {
		logger.Println("Error unmarshaling JSON:", err)
		return nil
	}
	return []byte(jsonBody.URL)
}

type JSONItem struct {
	Hash string `json:"short_url"`
	URL  string `json:"original_url"`
}

func urlsJSON(memData map[string]string, baseURL string) []byte {
	var jsonData []JSONItem
	for k, v := range memData {
		item := JSONItem{Hash: baseURL + "/" + k, URL: v}
		jsonData = append(jsonData, item)
	}
	jsonBytes, err := json.Marshal(jsonData)

	if err != nil {
		logger.Infoln(err)
	}

	return jsonBytes
}
