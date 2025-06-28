package usecases

import (
	"context"

	"github.com/aube/url-shortener/internal/app/hasher"
	"github.com/aube/url-shortener/internal/app/store"
	"github.com/aube/url-shortener/internal/logger"
)

func SaveURL(ctx context.Context, originalURL []byte, baseURL string) (string, error) {

	store := store.NewStore()
	log := logger.WithContext(ctx)

	hash := hasher.CalcHash(originalURL)
	err := store.Set(ctx, hash, string(originalURL))

	if err != nil {
		log.Error("SaveURL", "err", err)
		return "", err
	}

	shortURL := baseURL + "/" + hash

	log.Info("SaveURL", "hash", hash)

	return shortURL, nil
}
