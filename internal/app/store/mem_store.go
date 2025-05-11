package store

import (
	"context"
	"fmt"

	appErrors "github.com/aube/url-shortener/internal/app/apperrors"
	"github.com/aube/url-shortener/internal/logger"
)

type MemoryStore struct {
	s map[string]string
}

func (s *MemoryStore) Get(ctx context.Context, key string) (value string, ok bool) {
	log := logger.WithContext(ctx)

	value, ok = s.s[key]
	log.Info("Get", "key", key, "value", value)
	return value, ok
}

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

func (s *MemoryStore) Ping(ctx context.Context) error {
	return nil
}

func (s *MemoryStore) List(ctx context.Context) (map[string]string, error) {
	return s.s, nil
}

func (s *MemoryStore) SetMultiple(ctx context.Context, items map[string]string) error {
	log := logger.WithContext(ctx)

	for key, value := range items {
		log.Info("Set", "key", key, "value", value)
		s.s[key] = value
	}
	return nil
}

func (s *MemoryStore) Delete(ctx context.Context, hashes []string) error {
	log := logger.WithContext(ctx)

	for _, v := range hashes {
		log.Info("Delete", "hash", v)
		s.s[v] = ""
	}
	return nil
}

func NewMemStore() Storage {
	return &MemoryStore{
		s: make(map[string]string),
	}
}
