package usecases

import (
	"context"

	"github.com/aube/url-shortener/internal/app/store"
	"github.com/aube/url-shortener/internal/logger"
)

// StoragePing defines the interface for checking storage connectivity.
type StoragePing interface {
	// Ping checks if the storage is available.
	// Returns an error if the storage is not reachable.
	Ping(ctx context.Context) error
}

// StorePing checks the connectivity to the storage backend.
// Returns an error if the storage is not reachable.
func StorePing(ctx context.Context) error {
	store := store.NewStore()
	log := logger.WithContext(ctx)

	err := store.Ping(ctx)

	if err != nil {
		log.Debug("StorePing", "err", err)
		return err
	}

	return nil
}
