package store

import (
	"os"
)

var path string

func NewFileStore(storagePath string) {
	if err := os.MkdirAll(storagePath, os.ModePerm); err != nil {
		panic(err)
	}
	path = storagePath
}

func ReadFile(key string) (value string, ok bool) {
	content, err := os.ReadFile(path + "/" + key)
	return string(content), err == nil
}

func WriteFile(key string, value string) error {
	err := os.WriteFile(path+"/"+key, []byte(value), 0644)
	return err
}
