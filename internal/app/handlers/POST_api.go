package handlers

import (
	"io"
	"net/http"

	"github.com/aube/url-shortener/internal/app/usecases"
	"github.com/aube/url-shortener/internal/logger"
)

var usecasesSaveURL = usecases.SaveURL

// HandlerAPI create short URL in JSON
// @Summary Shorten a URL
// @Description Creates a short URL from a provided original URL
// @Tags URLs
// @Accept json
// @Produce json
// @Param request body string true "URL to shorten" example:"https://example.com"
// @Success 201 {object} map[string]string "URL created"
// @Success 409 {object} map[string]string "URL already exists"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/shorten [post]
func HandlerAPI(baseURL string) http.HandlerFunc {
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

		originalURL := readURLFromJSON(body)

		w.Header().Set("Content-Type", "application/json")
		httpStatus := http.StatusCreated

		shortURL, err := usecasesSaveURL(ctx, originalURL, baseURL)
		shortURL = `{"result":"` + shortURL + `"}`

		if err != nil {
			httpStatus = http.StatusConflict
		}
		w.WriteHeader(httpStatus)

		n, err := w.Write([]byte(shortURL))
		if err != nil {
			// Handle error (connection may have been closed)
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
			return
		}

		log.Info("HandlerAPI", "Wrote bytes", n)
	}
}
