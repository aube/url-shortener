package app

import (
	"context"
	"net/http"
	"time"

	"github.com/aube/url-shortener/internal/app/config"
	"github.com/aube/url-shortener/internal/app/router"
	"github.com/aube/url-shortener/internal/app/store"
	"github.com/aube/url-shortener/internal/logger"
)

type StorageGet interface {
	Get(ctx context.Context, key string) (value string, ok bool)
}
type StorageList interface {
	List(ctx context.Context) map[string]string
}
type StoragePing interface {
	Ping() error
}
type StorageSet interface {
	Set(ctx context.Context, key string, value string) error
}
type StorageSetMultiple interface {
	SetMultiple(ctx context.Context, l map[string]string) error
}
type Storage interface {
	StorageGet
	StorageList
	StoragePing
	StorageSet
	StorageSetMultiple
}

func Run() error {
	config := config.NewConfig()

	var storage Storage

	if config.DatabaseDSN != "" {
		storage = store.NewDBStore(config.DatabaseDSN)
	} else if config.FileStoragePath != "" {
		storage = store.NewFileStore(config.FileStoragePath)
	} else {
		storage = store.NewMemStore()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	router := router.Connect(ctx, storage)

	address := config.ServerHost + ":" + config.ServerPort
	logger.Infoln("Server starting at", address)

	err := http.ListenAndServe(address, router)

	if err != nil {
		logger.Infoln("Error starting server:", err)
		return err
	}

	return nil
}
