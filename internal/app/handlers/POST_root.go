package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	appErrors "github.com/aube/url-shortener/internal/app/apperrors"
	"github.com/aube/url-shortener/internal/app/hasher"
	"github.com/aube/url-shortener/internal/logger"
)

// HandlerRoot generate short URL
// @Summary Shorten a URL (text/plain)
// @Description Creates a short URL from a provided original URL (plain text input)
// @Tags URLs
// @Accept text/plain
// @Produce text/plain
// @Param request body string true "URL to shorten" example:"https://example.com"
// @Success 201 {string} string "Short URL"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router / [post]
func HandlerRoot(store StorageSet, baseURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logger.WithContext(ctx)

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

		originalURL := body
		contentType := r.Header.Get("Content-Type")
		responseContentType := contentType
		contentTypeJSON := strings.Contains(contentType, "application/json")
		acceptHeaderJSON := strings.Contains(r.Header.Get("Accept"), "application/json")

		responseContentJSON := contentTypeJSON || acceptHeaderJSON

		log.Debug(
			"HandlerRoot",
			"Request contentType:", contentType,
			"Response contentType:", responseContentType,
		)

		if responseContentJSON {
			originalURL = readURLFromJSON(body)
			responseContentType = "application/json"
		}
		w.Header().Set("Content-Type", responseContentType)

		hash := hasher.CalcHash(originalURL)
		httpStatus := http.StatusCreated

		err = store.Set(r.Context(), hash, string(originalURL))

		var herr *appErrors.HTTPError
		if errors.As(err, &herr) {
			httpStatus = herr.Code
		}

		w.WriteHeader(httpStatus)

		shortURL := baseURL + "/" + hash
		if responseContentJSON {
			fmt.Fprintf(w, `{"result":"%s"}`, shortURL)
		} else {
			fmt.Fprintf(w, "%s", shortURL)
		}

		log.Debug("URL:", shortURL, httpStatus)
	}
}
