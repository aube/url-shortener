package store

import (
	"context"
	"fmt"

	appErrors "github.com/aube/url-shortener/internal/app/apperrors"
	"github.com/aube/url-shortener/internal/logger"
)

type StorageGet interface {
	Get(ctx context.Context, key string) (value string, ok bool)
}
type StorageList interface {
	List(ctx context.Context) (map[string]string, error)
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
type StorageDelete interface {
	Delete(ctx context.Context, l []interface{}) error
}
type MemStorage interface {
	StorageGet
	StorageList
	StoragePing
	StorageSet
	StorageSetMultiple
	StorageDelete
}
type MemoryStore struct {
	s map[string]string
}

func (s *MemoryStore) Get(ctx context.Context, key string) (value string, ok bool) {
	value, ok = s.s[key]
	logger.Infoln("Get key:", key, value)
	return value, ok
}

func (s *MemoryStore) Set(ctx context.Context, key string, value string) error {
	if key == "" || value == "" {
		return fmt.Errorf("invalid input")
	}

	if _, ok := s.s[key]; ok {
		return appErrors.NewHTTPError(409, "conflict")
	}

	logger.Infoln("Set key:", key, value)
	s.s[key] = value

	return nil
}

func (s *MemoryStore) Ping() error {
	return nil
}

func (s *MemoryStore) List(ctx context.Context) (map[string]string, error) {
	return s.s, nil
}

func (s *MemoryStore) SetMultiple(ctx context.Context, items map[string]string) error {
	for k, v := range items {
		logger.Infoln("Set key:", k, v)
		s.s[k] = v
	}
	return nil
}

func (s *MemoryStore) Delete(ctx context.Context, hashes []interface{}) error {
	for _, v := range hashes {
		logger.Infoln("Del hash:", v)
		s.s[v.(string)] = ""
	}
	return nil
}

func NewMemStore() MemStorage {
	return &MemoryStore{
		s: make(map[string]string),
	}
}
