package store

import (
	"sync"

	"github.com/aube/url-shortener/internal/app/config"
)

var (
	storage Storage
	once    sync.Once
)

// MewStore creates and returns the appropriate storage implementation
// based on the application configuration. It checks the configuration
// in this order:
//  1. If DatabaseDSN is configured, returns a PostgreSQL DBStore
//  2. If FileStoragePath is configured, returns a FileStore
//  3. Otherwise returns an in-memory MemoryStore
//
// This function serves as a singleton for the storage implementations.
func NewStore() Storage {
	once.Do(func() {
		config := config.NewConfig()

		if config.DatabaseDSN != "" {
			storage = NewDBStore(config.DatabaseDSN)
		} else if config.FileStoragePath != "" {
			storage = NewFileStore(config.FileStoragePath)
		}
		storage = NewMemStore()
	})
	return storage
}

// Close required for propper termination of the connection to storage
func Close() error {
	config := config.NewConfig()
	if config.DatabaseDSN != "" {
		return CloseDBStore()
	}
	return nil
}
