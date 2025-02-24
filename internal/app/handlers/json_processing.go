package handlers

import (
	"encoding/json"
	"fmt"
)

type URLJson struct {
	URL string `json:"URL"`
}

func readURLFromJSON(body []byte) []byte {
	var jsonBody URLJson

	err := json.Unmarshal(body, &jsonBody)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return nil
	}
	return []byte(jsonBody.URL)
}
