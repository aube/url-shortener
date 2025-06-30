package restapi

import (
	"context"
	"net/http"

	"github.com/aube/url-shortener/internal/app/usecases"
	"github.com/aube/url-shortener/internal/logger"
)

var usecasesStorePing = usecases.StorePing

// StoragePing interface
type StoragePing interface {
	Ping(ctx context.Context) error
}

// HandlerPing ping database
// @Summary Check database connection
// @Description Verifies if the application can connect to the database
// @Tags Health
// @Produce text/plain
// @Success 200 {string} string "pong"
// @Failure 400 {object} map[string]string "Connection failed"
// @Router /ping [get]
func HandlerPing() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logger.WithContext(ctx)

		err := usecasesStorePing(ctx)

		if err != nil {
			log.Debug("HandlerPing", "err", err)
			http.Error(w, "URL not found", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		n, err := w.Write([]byte("pong"))
		if err != nil {
			// Handle error (connection may have been closed)
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
			return
		}

		log.Info("HandlerPing", "Wrote bytes", n)
	}
}
