package usecases

import (
	"context"
	"encoding/json"

	"github.com/aube/url-shortener/internal/app/store"
	"github.com/aube/url-shortener/internal/logger"
)

// StorageDelete interface
type StorageDelete interface {
	Delete(c context.Context, l []string) error
}

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
