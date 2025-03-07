package store

import (
	"fmt"

	"github.com/aube/url-shortener/internal/logger"
)

type Storage interface {
	Get(key string) (value string, ok bool)
	Set(key string, value string) error
	List() map[string]string
}

type MemoryStore struct {
	s map[string]string
}

var memData *MemoryStore

func init() {
	// initialize the singleton object. Here we're using a simple map as our storage
	memData = &MemoryStore{make(map[string]string)}
}

func NewMemStore() *MemoryStore {
	return memData
}

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

	WriteToFile(key, value)

	return nil
}

func (s *MemoryStore) List() map[string]string {
	return memData.s
}
