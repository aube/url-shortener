package store

import (
	"errors"
	"fmt"

	"github.com/aube/url-shortener/internal/logger"
)

type Storage interface {
	Get(key string) (value string, ok bool)
	List() map[string]string
	Ping() error

	Set(key string, value string) error
	SetMultiple(map[string]string) error
}

type MemoryStore struct {
	s map[string]string
}

// ErrConflict указывает на конфликт данных в хранилище.
var ErrConflict = errors.New("data conflict")

var memData = &MemoryStore{s: make(map[string]string)}

func (s *MemoryStore) Get(key string) (value string, ok bool) {
	value, ok = memData.s[key]
	logger.Infoln("Get key:", key, value)
	return value, ok
}

func (s *MemoryStore) Set(key string, value string) error {
	if key == "" || value == "" {
		return fmt.Errorf("invalid input")
	}

	if _, ok := memData.s[key]; ok {
		return ErrConflict
	}

	logger.Infoln("Set key:", key, value)
	memData.s[key] = value

	return nil
}

func (s *MemoryStore) Ping() error {
	return nil
}

func (s *MemoryStore) List() map[string]string {
	return memData.s
}

func (s *MemoryStore) SetMultiple(items map[string]string) error {
	for k, v := range items {
		logger.Infoln("Set key:", k, v)
		memData.s[k] = v
	}
	return nil
}

func NewMemStore() Storage {
	return &MemoryStore{}
}
