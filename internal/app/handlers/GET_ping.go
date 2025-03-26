package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aube/url-shortener/internal/logger"
)

type StoragePing interface {
	Ping(ctx context.Context) error
}

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
