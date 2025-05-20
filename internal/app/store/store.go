package store

import "github.com/aube/url-shortener/internal/app/config"

// MewStore creates and returns the appropriate storage implementation
// based on the application configuration. It checks the configuration
// in this order:
//  1. If DatabaseDSN is configured, returns a PostgreSQL DBStore
//  2. If FileStoragePath is configured, returns a FileStore
//  3. Otherwise returns an in-memory MemoryStore
//
// This function serves as a factory for the storage implementations.
func MewStore() Storage {
	config := config.NewConfig()
	if config.DatabaseDSN != "" {
		return NewDBStore(config.DatabaseDSN)
	} else if config.FileStoragePath != "" {
		return NewFileStore(config.FileStoragePath)
	}
	return NewMemStore()
}
