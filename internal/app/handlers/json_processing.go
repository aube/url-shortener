package handlers

import (
	"encoding/json"

	"github.com/aube/url-shortener/internal/logger"
)

// URLJson JSON struct
type URLJson struct {
	URL string `json:"URL"`
}

func readURLFromJSON(body []byte) []byte {
	log := logger.Get()

	var jsonBody URLJson

	err := json.Unmarshal(body, &jsonBody)
	if err != nil {
		log.Debug("readURLFromJSON", "err", err)
		return nil
	}
	return []byte(jsonBody.URL)
}
