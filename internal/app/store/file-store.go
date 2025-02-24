package store

import (
	"fmt"
	"os"
	"strings"
)

var file string

func getDirFromPath(path string) (dir string) {
	parts := strings.Split(path, `/`)
	return strings.Join(parts[:len(parts)-1], "/")
}

func NewFileStore(storagePath string) {
	d := getDirFromPath(storagePath)

	fmt.Println("create dir:", d)

	if err := os.MkdirAll(d, os.ModePerm); err != nil {
		panic(err)
	}

	data := []byte("{}")

	f, err := os.Create(storagePath)

	fmt.Println("create dir:", storagePath)

	if err != nil {
		fmt.Println("Unable to create file:", err)
		panic(err)
	}

	defer f.Close()
	f.Write(data)

	file = storagePath
}

func ReadFile(key string) (value string, ok bool) {
	content, err := os.ReadFile(file)
	return string(content), err == nil
}

func WriteFile(key string, value string) error {
	err := os.WriteFile(file, []byte(value), 0644)
	return err
}
