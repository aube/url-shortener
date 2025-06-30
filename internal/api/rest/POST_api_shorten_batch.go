package restapi

import (
	"io"
	"net/http"

	"github.com/aube/url-shortener/internal/app/usecases"
	"github.com/aube/url-shortener/internal/logger"
)

var usecasesSaveURLS = usecases.SaveURLS

// HandlerShortenBatch create multiple short URLs
// @Summary Shorten multiple URLs
// @Description Creates short URLs for multiple provided original URLs
// @Tags URLs
// @Accept json
// @Produce json
// @Param request body []handlers.inputBatchJSONItem true "Batch of URLs to shorten"
// @Success 201 {array} handlers.outputBatchJSONItem
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/shorten/batch [post]
func HandlerShortenBatch(baseURL string) http.HandlerFunc {
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

		outputJSON, err := usecasesSaveURLS(ctx, body, baseURL)

		if err != nil {
			log.Error("SetMultiple", "err", err)
			http.Error(w, "Failed to write URLs", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		n, err := w.Write(outputJSON)
		if err != nil {
			// Handle error (connection may have been closed)
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
			return
		}

		log.Info("HandlerShortenBatch", "Wrote bytes", n)
	}
}
