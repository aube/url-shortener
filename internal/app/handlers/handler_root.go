package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/aube/url-shortener/internal/app/hasher"
	"github.com/aube/url-shortener/internal/logger"
)

func HandlerRoot(MemoryStore StorageSet, baseURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		responseContentType := contentType
		contentTypeJSON := strings.Contains(contentType, "application/json")
		acceptHeaderJSON := strings.Contains(r.Header.Get("Accept"), "application/json")

		responseContentJSON := contentTypeJSON || acceptHeaderJSON

		logger.Println(
			"Request contentType:", contentType,
			"Response contentType:", responseContentType,
		)

		if responseContentJSON {
			originalURL = readURLFromJSON(body)
			responseContentType = "application/json"
		}

		hash := hasher.CalcHash(originalURL)
		MemoryStore.Set(hash, string(originalURL))

		shortURL := baseURL + "/" + hash

		w.Header().Set("Content-Type", responseContentType)
		w.WriteHeader(http.StatusCreated)

		if responseContentJSON {
			fmt.Fprintf(w, `{"result":"%s"}`, shortURL)
		} else {
			fmt.Fprintf(w, "%s", shortURL)
		}

		logger.Println("URL:", shortURL, http.StatusCreated)
	}
}
