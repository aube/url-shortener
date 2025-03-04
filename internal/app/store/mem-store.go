package store

import (
	"encoding/json"
	"fmt"

	"github.com/aube/url-shortener/internal/logger"
)

type Storage interface {
	Get(key string) (value string, ok bool)
	Set(key string, value string) error
}

type MemoryStore struct {
	s map[string]string
}

var memData *MemoryStore

func init() {
	// initialize the singleton object. Here we're using a simple map as our storage
	memData = &MemoryStore{make(map[string]string)}
}

func NewMemoryStore() *MemoryStore {
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

type JSONItem struct {
	Hash string `json:"short_url"`
	URL  string `json:"original_url"`
}

func (s *MemoryStore) JSON(baseURL string) []byte {
	var jsonData []JSONItem
	for k, v := range memData.s {
		item := JSONItem{Hash: baseURL + "/" + k, URL: v}
		jsonData = append(jsonData, item)
	}
	jsonBytes, err := json.Marshal(jsonData)

	if err != nil {
		logger.Infoln(err)
	}

	return jsonBytes
}

func SetValue(key string, value string) {
	memData.s[key] = value
}
