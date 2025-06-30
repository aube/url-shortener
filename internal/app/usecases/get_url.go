package usecases

import (
	"context"
	"errors"

	"github.com/aube/url-shortener/internal/app/store"
	"github.com/aube/url-shortener/internal/logger"
)

// StorageGet defines the interface for retrieving a single URL from storage.
type StorageGet interface {
	// Get retrieves the original URL associated with the given hash.
	// Returns the URL and a boolean indicating if it was found.
	Get(c context.Context, key string) (value string, ok bool)
}

// GetURL retrieves the original URL associated with the given hash.
// Returns the URL or an error if not found.
func GetURL(ctx context.Context, id string) (string, error) {
	store := store.NewStore()
	log := logger.WithContext(ctx)

	url, ok := store.Get(ctx, id)

	log.Debug("GetURL", "id", id, "url", url, "ok", ok)

	if !ok {
		return "", errors.New("not found")
	}

	return url, nil
}
