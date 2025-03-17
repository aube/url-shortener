package store

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/aube/url-shortener/internal/logger"
)

type FileStorage interface {
	StorageGet
	StorageList
	StoragePing
	StorageSet
	StorageSetMultiple
}

type FileStore struct {
	s          map[string]string
	pathToFile string
}

func (s *FileStore) Get(ctx context.Context, key string) (value string, ok bool) {
	value, ok = s.s[key]
	logger.Infoln("Get key:", key, value)
	return value, ok
}

func (s *FileStore) Set(ctx context.Context, key string, value string) error {
	if key == "" || value == "" {
		return fmt.Errorf("invalid input")
	}

	if _, ok := s.s[key]; ok {
		return ErrConflict
	}

	logger.Infoln("Set key:", key, value)
	s.s[key] = value

	WriteToFile(key, value, s.pathToFile)

	return nil
}

func (s *FileStore) List(ctx context.Context) map[string]string {
	return s.s
}

func (s *FileStore) Ping() error {
	return nil
}

func (s *FileStore) SetMultiple(ctx context.Context, items map[string]string) error {
	for k, v := range items {
		logger.Infoln("Set key:", k, v)
		s.s[k] = v

		WriteToFile(k, v, s.pathToFile)
	}
	return nil
}

func getDirFromPath(path string) (dir string) {
	parts := strings.Split(path, `/`)
	return strings.Join(parts[:len(parts)-1], "/")
}

func createDir(storagePath string) {
	d := getDirFromPath(storagePath)

	logger.Println("create dir:", d)

	if err := os.MkdirAll(d, os.ModePerm); err != nil {
		panic(err)
	}
}

func createFile(storagePath string) {
	if _, err := os.Stat(storagePath); err == nil {
		// file exists
		return
	}

	data := []byte("")
	f, err := os.Create(storagePath)
	logger.Println("create file:", storagePath)

	if err != nil {
		logger.Println("Unable to create file:", err)
		panic(err)
	}
	defer f.Close()
	f.Write(data)
}

type itemURL struct {
	Hash string `json:"Hash"`
	URL  string `json:"OriginalURL"`
}

func lineToJSON(line string) itemURL {
	req := itemURL{}
	if err := json.Unmarshal([]byte(line), &req); err != nil {
		panic(err)
	}
	return req
}

func getFileContent(storagePath string) map[string]string {
	file, err := os.Open(storagePath)
	if err != nil {
		logger.Println(err)
	}
	defer file.Close()

	data := make(map[string]string)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			json := lineToJSON(line)
			data[json.Hash] = json.URL
		}
	}

	if err = scanner.Err(); err != nil {
		logger.Println(err)
	}

	return data
}

func NewFileStore(storagePath string) FileStorage {
	createDir(storagePath)
	createFile(storagePath)
	data := getFileContent(storagePath)

	return &FileStore{
		pathToFile: storagePath,
		s:          data,
	}
}

func WriteToFile(key string, value string, pathToFile string) error {
	f, err := os.OpenFile(pathToFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	json, err := json.Marshal(itemURL{Hash: key, URL: value})
	if err != nil {
		return err
	}

	if _, err = f.WriteString(string(json) + "\n"); err != nil {
		return err
	}

	logger.Println("WriteToFile:", json)
	return nil
}
