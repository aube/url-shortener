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
