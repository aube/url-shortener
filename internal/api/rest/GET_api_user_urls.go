package restapi

import (
	"encoding/json"
	"errors"
	"net/http"

	appErrors "github.com/aube/url-shortener/internal/app/apperrors"
	"github.com/aube/url-shortener/internal/app/usecases"
	"github.com/aube/url-shortener/internal/logger"
)

// HandlerAPIUserUrls read multiple URLs for a user
// @Summary Get user URLs
// @Description Returns all shortened URLs belonging to the authenticated user
// @Tags URLs
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} handlers.JSONItem
// @Success 204 "No content"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/user/urls [get]
func HandlerAPIUserUrls(baseURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logger.WithContext(ctx)
		w.Header().Set("Content-Type", "application/json")

		json, err := usecases.GetURLS(ctx, baseURL)

		var herr *appErrors.HTTPError
		if errors.As(err, &herr) {
			w.WriteHeader(herr.Code)
			return
		}

		if len(json) == 2 { // "[]"
			w.WriteHeader(204)
			return
		}

		if err != nil {
			log.Error("getJSON", "err", err)
		}

		w.WriteHeader(http.StatusOK)
		n, err := w.Write(json)

		if err != nil {
			// Handle error (connection may have been closed)
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
			return
		}

		log.Info("getJSON", "Wrote bytes", n)
	}
}

// JSONItem struct
type JSONItem struct {
	Hash string `json:"short_url"`
	URL  string `json:"original_url"`
}

func getJSON(memData map[string]string, baseURL string) ([]byte, error) {

	var jsonData []JSONItem

	for k, v := range memData {
		item := JSONItem{Hash: baseURL + "/" + k, URL: v}
		jsonData = append(jsonData, item)
	}

	return json.Marshal(jsonData)
}
