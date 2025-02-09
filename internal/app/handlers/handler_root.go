package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aube/url-shortener/internal/app/hashes"
)

type URLJson struct {
	URL string `json:"URL"`
}

func readUrlFromJson(body []byte) []byte {
	var jsonBody URLJson

	err := json.Unmarshal(body, &jsonBody)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		panic("ololo")
	}
	return []byte(jsonBody.URL)
}

func HandlerRoot(w http.ResponseWriter, r *http.Request, baseUrl string) {
	switch r.Method {
	case "POST":
		// Read the entire body content
		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			break
		}
		defer r.Body.Close()

		if len(body) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			break
		}
		originalUrl := body
		contentType := r.Header.Get("Content-Type")

		fmt.Println("contentType:", contentType)

		if contentType == "application/json" {
			originalUrl = readUrlFromJson(body)
		}

		hash := hashes.SetURLHash(originalUrl)
		shortUrl := baseUrl + "/" + hash

		if contentType == "application/json" {
			w.Header().Set("Content-Type", "application/json")
			shortUrl = `{"result":"` + shortUrl + `"}`
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, shortUrl)
		fmt.Println("URL:", shortUrl, http.StatusCreated)
	default:
		fmt.Println("Not served method:", r.Method)
	}
}
