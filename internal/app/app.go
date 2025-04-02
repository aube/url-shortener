package app

import (
	"context"
	"log"
	"net/http"

	"github.com/aube/url-shortener/internal/app/config"
	"github.com/aube/url-shortener/internal/app/router"
	"github.com/aube/url-shortener/internal/app/store"
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
	Delete(ctx context.Context, l []string) error
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

	router := router.Connect(storage)

	address := config.ServerHost + ":" + config.ServerPort
	log.Println("Server starting", "address", address)

	err := http.ListenAndServe(address, router)

	if err != nil {
		log.Fatal("Starting server", "err", err)
	}

	return nil
}
