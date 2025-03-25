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
	List(ctx context.Context) (map[string]string, error)
}
type StoragePing interface {
	Ping(ctx context.Context) error
}
type StorageSet interface {
	Set(ctx context.Context, key string, value string) error
}
type StorageSetMultiple interface {
	SetMultiple(ctx context.Context, l map[string]string) error
}
type StorageDelete interface {
	Delete(ctx context.Context, l []interface{}) error
}

type Storage interface {
	StorageGet
	StorageList
	StoragePing
	StorageSet
	StorageSetMultiple
	StorageDelete
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

	// TODO разобраться с контекстом, пока ставлю таймаут побольше
	ctx, cancel := context.WithTimeout(context.Background(), 5000*time.Second)
	defer cancel()

	log := logger.WithContext(ctx)

	router := router.Connect(ctx, storage)

	address := config.ServerHost + ":" + config.ServerPort
	log.Info("Server starting", "address", address)

	err := http.ListenAndServe(address, router)

	if err != nil {
		log.Error("Starting server", "err", err)
		panic(err)
	}

	return nil
}
