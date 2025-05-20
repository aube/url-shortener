package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aube/url-shortener/internal/logger"
)

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
func HandlerPing(store StoragePing) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logger.WithContext(ctx)

		err := store.Ping(ctx)

		if err != nil {
			log.Debug("HandlerPing", "err", err)
			http.Error(w, "URL not found", http.StatusBadRequest)
			return
		}

		log.Debug("Ping DB")
		fmt.Fprintf(w, `pong`)
		w.WriteHeader(http.StatusOK)
	}
}
