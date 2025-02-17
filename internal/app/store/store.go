package store

import (
	"fmt"
)

type Storer interface {
	Get(key string) (value string, ok bool)
	Set(key string, value string) error
}

type MemoryStore struct {
	data map[string]string
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{data: make(map[string]string)}
}

func (s *MemoryStore) Get(key string) (value string, ok bool) {
	value, ok = s.data[key]
	return value, ok
}

func (s *MemoryStore) Set(key string, value string) error {
	if key == "" || value == "" {
		return fmt.Errorf("invalid input")
	}
	s.data[key] = value
	return nil
}
