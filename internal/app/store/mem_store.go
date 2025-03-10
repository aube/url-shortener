package store

import (
	"fmt"

	"github.com/aube/url-shortener/internal/logger"
)

type Storage interface {
	Get(key string) (value string, ok bool)
	Set(key string, value string) error
	List() map[string]string
	Ping() error
}

type MemoryStore struct {
	s map[string]string
}

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

func NewMemStore() Storage {
	return &MemoryStore{}
}
