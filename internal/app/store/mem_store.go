package store

import (
	"context"
	"fmt"
	"time"

	appErrors "github.com/aube/url-shortener/internal/app/apperrors"
	"github.com/aube/url-shortener/internal/logger"
)

// MemoryStore is an in-memory implementation of the Storage interface
// using a map to store URL mappings. Suitable for development and testing.
type MemoryStore struct {
	s map[string]string
}

// Get retrieves a URL by its shortened key from memory.
// Returns the URL and true if found, empty string and false otherwise.
func (s *MemoryStore) Get(ctx context.Context, key string) (value string, ok bool) {
	log := logger.WithContext(ctx)

	value, ok = s.s[key]
	log.Info("Get", "key", key, "value", value)
	return value, ok
}

// Set stores a new URL mapping in memory.
// Returns an error if the key is empty, value is empty, or if the key already exists.
func (s *MemoryStore) Set(ctx context.Context, key string, value string) error {
	log := logger.WithContext(ctx)

	if key == "" || value == "" {
		return fmt.Errorf("invalid input")
	}

	if _, ok := s.s[key]; ok {
		return appErrors.NewHTTPError(409, "conflict")
	}

	log.Info("Set", "key", key, "value", value)
	s.s[key] = value

	return nil
}

// Ping always returns nil for in-memory storage as it's always available.
func (s *MemoryStore) Ping(ctx context.Context) error {
	return nil
}

// List returns all URL mappings currently stored in memory.
func (s *MemoryStore) List(ctx context.Context) (map[string]string, error) {
	return s.s, nil
}

// SetMultiple stores multiple URL mappings in a batch operation.
func (s *MemoryStore) SetMultiple(ctx context.Context, items map[string]string) error {
	log := logger.WithContext(ctx)

	for key, value := range items {
		log.Info("Set", "key", key, "value", value)
		s.s[key] = value
	}
	return nil
}

// Delete marks one or more URLs as deleted by setting their values to empty string.
func (s *MemoryStore) Delete(ctx context.Context, hashes []string) error {
	log := logger.WithContext(ctx)

	for _, v := range hashes {
		log.Info("Delete", "hash", v)
		s.s[v] = ""
	}
	return nil
}

// Stats select amount of urls and users from database.
func (s *MemoryStore) Stats(ctx context.Context) (urls int, users int, err error) {
	// log := logger.WithContext(ctx)

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// row := db.QueryRowContext(ctx, postgre.selectDBStatistics, key)
	// var originalURL string
	// var deleted bool

	// err := row.Scan(&originalURL, &deleted)

	// if err != nil {
	// 	log.Error("Get", "err", err)
	// 	return "", false
	// }
	// if deleted {
	// 	return "", true
	// }

	return 1, 1, nil
}

// NewMemStore creates and returns a new in-memory storage instance.
func NewMemStore() Storage {
	return &MemoryStore{
		s: make(map[string]string),
	}
}
