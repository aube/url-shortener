package store

import (
	"context"
	"fmt"

	appErrors "github.com/aube/url-shortener/internal/app/apperrors"
	"github.com/aube/url-shortener/internal/app/ctxkeys"
	"github.com/aube/url-shortener/internal/logger"
)

// MemoryStore is an in-memory implementation of the Storage interface
// using a map to store URL mappings. Suitable for development and testing.
type MemoryStore struct {
	urls  map[string]string
	users map[string]string
}

// Get retrieves a URL by its shortened key from memory.
// Returns the URL and true if found, empty string and false otherwise.
func (s *MemoryStore) Get(ctx context.Context, key string) (value string, ok bool) {
	log := logger.WithContext(ctx)
	userID := ctx.Value(ctxkeys.UserIDKey).(string)

	value, ok = s.urls[key]
	if ok && s.users[key] == userID {
		log.Info("Get", "key", key, "value", value)
		return value, ok
	}
	return "", false
}

// Set stores a new URL mapping in memory.
// Returns an error if the key is empty, value is empty, or if the key already exists.
func (s *MemoryStore) Set(ctx context.Context, key string, value string) error {
	log := logger.WithContext(ctx)
	userID := ctx.Value(ctxkeys.UserIDKey).(string)

	if key == "" || value == "" {
		return fmt.Errorf("invalid input")
	}

	if _, ok := s.urls[key]; ok {
		return appErrors.NewHTTPError(409, "conflict")
	}

	log.Info("Set", "key", key, "value", value)
	s.urls[key] = value
	s.users[key] = userID

	return nil
}

// Ping always returns nil for in-memory storage as it's always available.
func (s *MemoryStore) Ping(ctx context.Context) error {
	return nil
}

// List returns all URL mappings currently stored in memory.
func (s *MemoryStore) List(ctx context.Context) (map[string]string, error) {
	userID := ctx.Value(ctxkeys.UserIDKey).(string)
	urls := make(map[string]string)

	for hash, url := range s.urls {
		if s.users[hash] == userID {
			urls[hash] = url
		}
	}
	return urls, nil
}

// SetMultiple stores multiple URL mappings in a batch operation.
func (s *MemoryStore) SetMultiple(ctx context.Context, items map[string]string) error {
	log := logger.WithContext(ctx)
	userID := ctx.Value(ctxkeys.UserIDKey).(string)

	for key, value := range items {
		log.Info("Set", "key", key, "value", value)
		s.urls[key] = value
		s.users[key] = userID
	}
	return nil
}

// Delete marks one or more URLs as deleted by setting their values to empty string.
func (s *MemoryStore) Delete(ctx context.Context, hashes []string) error {
	log := logger.WithContext(ctx)
	userID := ctx.Value(ctxkeys.UserIDKey).(string)

	for _, hash := range hashes {
		if s.users[hash] == userID {
			log.Info("Delete", "hash", hash)
			delete(s.urls, hash)
		}
	}
	return nil
}

// Stats select amount of urls and users from database.
func (s *MemoryStore) Stats(ctx context.Context) (int, int, error) {
	log := logger.WithContext(ctx)
	urls := len(s.urls)
	users := make(map[string]string)

	for _, user := range s.users {
		users[user] = ""
	}

	log.Info("Stats", "urls", urls, "users", len(users))
	return urls, len(users), nil
}

// NewMemStore creates and returns a new in-memory storage instance.
func NewMemStore() Storage {
	return &MemoryStore{
		urls:  make(map[string]string),
		users: make(map[string]string),
	}
}
