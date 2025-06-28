package usecases

import (
	"context"

	"github.com/aube/url-shortener/internal/app/store"
	"github.com/aube/url-shortener/internal/logger"
)

// StoragePing interface
type StoragePing interface {
	Ping(ctx context.Context) error
}

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
