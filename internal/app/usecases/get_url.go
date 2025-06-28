package usecases

import (
	"context"
	"errors"

	"github.com/aube/url-shortener/internal/app/store"
	"github.com/aube/url-shortener/internal/logger"
)

// StorageGet interface
type StorageGet interface {
	Get(c context.Context, key string) (value string, ok bool)
}

func GetURL(ctx context.Context, id string) (string, error) {
	store := store.NewStore()
	log := logger.WithContext(ctx)

	url, ok := store.Get(ctx, id)

	if !ok {
		return "", errors.New("not found")
	}

	log.Debug("GetURL", "id", id, "url", url)

	return url, nil
}
