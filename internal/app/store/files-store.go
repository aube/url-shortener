package store

import (
	"os"
)

var path string

func NewFilesStore(storageDir string) {
	if err := os.MkdirAll(storageDir, os.ModePerm); err != nil {
		panic(err)
	}
	path = storageDir
}

func ReadFiles(key string) (value string, ok bool) {
	content, err := os.ReadFile(path + "/" + key)
	return string(content), err == nil
}

func WriteFiles(key string, value string) error {
	err := os.WriteFile(path+"/"+key, []byte(value), 0644)
	return err
}
