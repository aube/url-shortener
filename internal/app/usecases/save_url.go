package usecases

import (
	"context"

	"github.com/aube/url-shortener/internal/app/hasher"
	"github.com/aube/url-shortener/internal/app/store"
	"github.com/aube/url-shortener/internal/logger"
)

// SaveURL stores a new URL and returns its shortened version.
// The URL is hashed and stored with the provided baseURL.
// Returns the shortened URL or an error if the operation fails.
func SaveURL(ctx context.Context, originalURL []byte, baseURL string) (string, error) {

	store := store.NewStore()
	log := logger.WithContext(ctx)

	hash := hasher.CalcHash(originalURL)
	shortURL := baseURL + "/" + hash

	err := store.Set(ctx, hash, string(originalURL))
	log.Info("SaveURL", "hash", hash)

	return shortURL, err

}
