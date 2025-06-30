package usecases

import (
	"context"
	"encoding/json"

	"github.com/aube/url-shortener/internal/app/store"
	"github.com/aube/url-shortener/internal/logger"
)

// StorageStats defines the interface for retrieving statistics about stored URLs.
type StorageStats interface {
	// Stats returns the count of URLs and users in storage.
	// Returns an error if the operation fails.
	Stats(ctx context.Context) (urls int, users int, err error)
}

// GetStats retrieves statistics about stored URLs and users.
// Returns the statistics as JSON bytes or an error if the operation fails.
func GetStats(ctx context.Context) ([]byte, error) {
	store := store.NewStore()
	log := logger.WithContext(ctx)

	urls, users, err := store.Stats(ctx)

	if err != nil {
		return nil, err
	}

	// JSONItem struct
	type response struct {
		Urls  int `json:"urls"`
		Users int `json:"users"`
	}

	json, err := json.Marshal(response{Urls: urls, Users: users})

	if err != nil {
		log.Error("HandlerInternalStats", "json.Marshal", err)
		return nil, err
	}

	log.Info("GetStats", "urls", urls, "users", users)

	return json, err
}
