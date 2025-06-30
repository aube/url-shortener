package restapi

import (
	"io"
	"net/http"

	"github.com/aube/url-shortener/internal/app/usecases"
	"github.com/aube/url-shortener/internal/logger"
)

var usecasesDeleteURLS = usecases.DeleteURLS

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
func HandlerAPIUserUrlsDel(baseURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logger.WithContext(ctx)

		body, err := io.ReadAll(r.Body)

		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}

		err = usecasesDeleteURLS(ctx, body)

		if err != nil {
			log.Error("HandlerAPIUserUrlsDel", "err", err)
			http.Error(w, "Failed to delete", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)

		log.Debug("HandlerAPIUserUrlsDel", "ok", true)
	}
}
