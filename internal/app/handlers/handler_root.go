package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/aube/url-shortener/internal/app/hashes"
)

type URLJson struct {
	URL string `json:"URL"`
}

func readURLFromJSON(body []byte) []byte {
	var jsonBody URLJson

	err := json.Unmarshal(body, &jsonBody)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		panic("ololo")
	}
	return []byte(jsonBody.URL)
}

func HandlerRoot(w http.ResponseWriter, r *http.Request, baseURL string) {
	switch r.Method {
	case "POST":
		// Read the entire body content
		body, err := io.ReadAll(r.Body)

		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			break
		}
		defer r.Body.Close()

		if len(body) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			break
		}
		originalURL := body
		contentType := r.Header.Get("Content-Type")

		fmt.Println("contentType:", contentType)

		if contentType == "application/json" {
			originalURL = readURLFromJSON(body)
		}

		hash := hashes.SetURLHash(originalURL)
		shortURL := baseURL + "/" + hash

		w.WriteHeader(http.StatusCreated)
		if contentType == "application/json" {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"result":"%s"}`, shortURL)
		} else {
			fmt.Fprintf(w, "%s", shortURL)
		}

		fmt.Println("URL:", shortURL, http.StatusCreated)
	default:
		fmt.Println("Not served method:", r.Method)
	}
}
