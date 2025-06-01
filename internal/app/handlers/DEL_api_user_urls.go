package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

const numWorkers int = 10

// StorageDelete interface
type StorageDelete interface {
	Delete(c context.Context, l []string) error
}

// HandlerAPIUserUrlsDel deletes multiple URLs for a user
// @Summary Delete user URLs
// @Description Deletes multiple shortened URLs belonging to the authenticated user
// @Tags URLs
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body []string true "Array of URL hashes to delete"
// @Success 202 {string} string "Deletion request accepted"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/user/urls [delete]
func HandlerAPIUserUrlsDel(store StorageDelete, baseURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		body, err := io.ReadAll(r.Body)

		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}

		var data []string

		err = json.Unmarshal([]byte(body), &data)
		if err != nil {
			http.Error(w, "Failed to read JSON", http.StatusInternalServerError)
			return
		}

		err = store.Delete(r.Context(), data)
		if err != nil {
			http.Error(w, "Failed to delete", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}
