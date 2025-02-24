package store

import (
	"fmt"
)

// ???
/* type Storage interface {
	Get(key string) (value string, ok bool)
	Set(key string, value string) error
}
*/

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
	value, ok = ReadFiles(key)
	fmt.Println("Read key:", key, value, ok)
	if ok {
		return value, ok
	}

	value, ok = data.s[key]
	fmt.Println("Get key:", key, value)
	return value, ok
}

func (s *MemoryStore) Set(key string, value string) error {
	if key == "" || value == "" {
		return fmt.Errorf("invalid input")
	}

	err := WriteFiles(key, value)
	fmt.Println("Write key:", key, value)
	if err != nil {
		panic(err)
	}

	fmt.Println("Set key:", key, value)
	data.s[key] = value

	return nil
}
