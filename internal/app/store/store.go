package store

import (
	"fmt"
)

type Storage interface {
	Get(key string) (value string, ok bool)
	Set(key string, value string) error
}

type MemoryStore struct {
	s map[string]string
}

var data *MemoryStore

func init() {
	// initialize the singleton object. Here we're using a simple map as our storage
	data = &MemoryStore{make(map[string]string)}
}

func NewMemoryStore() *MemoryStore {
	return data
}
func (s *MemoryStore) Get(key string) (value string, ok bool) {
	value, ok = data.s[key]
	fmt.Println("Get key:value", key, value)
	return value, ok
}

func (s *MemoryStore) Set(key string, value string) error {
	if key == "" || value == "" {
		return fmt.Errorf("invalid input")
	}

	fmt.Println("Set key:value", key, value)
	data.s[key] = value

	return nil
}
