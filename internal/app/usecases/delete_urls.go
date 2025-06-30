package usecases

import (
	"context"
	"encoding/json"

	"github.com/aube/url-shortener/internal/app/store"
	"github.com/aube/url-shortener/internal/logger"
)

// StorageDelete defines the interface for deleting URLs from storage.
type StorageDelete interface {
	// Delete removes multiple URLs from storage based on their hashes.
	// Returns an error if the operation fails.
	Delete(c context.Context, l []string) error
}

// DeleteURLS handles the deletion of multiple URLs.
// It unmarshals the request body into a slice of URL hashes and deletes them from storage.
// Returns an error if the operation fails.
func DeleteURLS(ctx context.Context, body []byte) error {
	store := store.NewStore()
	log := logger.WithContext(ctx)

	var data []string

	err := json.Unmarshal(body, &data)

	if err != nil {
		log.Error("DeleteURLS", "err", err)
		return err
	}

	err = store.Delete(ctx, data)

	if err != nil {
		log.Error("DeleteURLS", "err2", err)
		return err
	}

	return nil
}
