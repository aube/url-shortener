package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/aube/url-shortener/internal/app/hasher"
	"github.com/aube/url-shortener/internal/app/store"
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

func HandlerRoot(baseURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		MemoryStore := store.NewMemoryStore()

		if r.Body == nil || r.ContentLength == 0 {
			http.Error(w, "Request body is empty", http.StatusBadRequest)
			return
		}

		// Read the entire body content
		body, err := io.ReadAll(r.Body)

		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		originalURL := body
		contentType := r.Header.Get("Content-Type")

		fmt.Println("contentType:", contentType)

		if contentType == "application/json" {
			originalURL = readURLFromJSON(body)
		}

		hash := hasher.CalcHash(originalURL)
		MemoryStore.Set(hash, string(originalURL))
		MemoryStore.Get(hash)

		shortURL := baseURL + "/" + hash

		w.WriteHeader(http.StatusCreated)
		if contentType == "application/json" {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"result":"%s"}`, shortURL)
		} else {
			fmt.Fprintf(w, "%s", shortURL)
		}

		fmt.Println("URL:", shortURL, http.StatusCreated)

	}
}
