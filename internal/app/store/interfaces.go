package store

import "context"

// StorageGet defines the interface for retrieving URL mappings.
type StorageGet interface {
	// Get retrieves the original URL for a given shortened key.
	// Returns the URL and a boolean indicating if the key was found.
	Get(ctx context.Context, key string) (value string, ok bool)
}

// StorageList defines the interface for listing URL mappings.
type StorageList interface {
	// List returns all URL mappings for the current user.
	// Returns a map of shortened keys to original URLs.
	List(ctx context.Context) (map[string]string, error)
}

// StoragePing defines the interface for checking storage availability.
type StoragePing interface {
	// Ping checks if the storage backend is available.
	Ping(ctx context.Context) error
}

// StorageSet defines the interface for storing URL mappings.
type StorageSet interface {
	// Set stores a new shortened URL mapping.
	// Returns an error if the key already exists (conflict).
	Set(ctx context.Context, key string, value string) error
}

// StorageSetMultiple defines the interface for batch URL storage operations.
type StorageSetMultiple interface {
	// SetMultiple stores multiple URL mappings in a single operation.
	SetMultiple(ctx context.Context, l map[string]string) error
}

// StorageDelete defines the interface for deleting URL mappings.
type StorageDelete interface {
	// Delete marks one or more URLs as deleted.
	Delete(ctx context.Context, l []string) error
}

// Storage is the comprehensive interface combining all storage operations.
type Storage interface {
	StorageGet
	StorageList
	StoragePing
	StorageSet
	StorageSetMultiple
	StorageDelete
}
