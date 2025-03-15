package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/aube/url-shortener/internal/logger"
)

type MemStorage interface {
	Get(ctx context.Context, key string) (value string, ok bool)
	List(ctx context.Context) map[string]string
	Ping() error
	Set(ctx context.Context, key string, value string) error
	SetMultiple(ctx context.Context, l map[string]string) error
}

type MemoryStore struct {
	s map[string]string
}

// ErrConflict указывает на конфликт данных в хранилище.
var ErrConflict = errors.New("data conflict")

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
		return ErrConflict
	}

	logger.Infoln("Set key:", key, value)
	s.s[key] = value

	return nil
}

func (s *MemoryStore) Ping() error {
	return nil
}

func (s *MemoryStore) List(ctx context.Context) map[string]string {
	return s.s
}

func (s *MemoryStore) SetMultiple(ctx context.Context, items map[string]string) error {
	for k, v := range items {
		logger.Infoln("Set key:", k, v)
		s.s[k] = v
	}
	return nil
}

func NewMemStore() MemStorage {
	return &MemoryStore{
		s: make(map[string]string),
	}
}
