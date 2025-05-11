package store

import "github.com/aube/url-shortener/internal/app/config"

func MewStore() Storage {

	config := config.NewConfig()
	if config.DatabaseDSN != "" {
		return NewDBStore(config.DatabaseDSN)
	} else if config.FileStoragePath != "" {
		return NewFileStore(config.FileStoragePath)
	}
	return NewMemStore()

}
