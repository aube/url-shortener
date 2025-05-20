package store

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	appErrors "github.com/aube/url-shortener/internal/app/apperrors"
	"github.com/aube/url-shortener/internal/logger"
)

// FileStore is a file-based implementation of the Storage interface
// that persists URL mappings to disk in JSON format.
type FileStore struct {
	s          map[string]string
	pathToFile string
}

// Get retrieves a URL by its shortened key from the file storage.
// Returns the URL and true if found, empty string and false otherwise.
func (s *FileStore) Get(ctx context.Context, key string) (value string, ok bool) {
	log := logger.WithContext(ctx)

	value, ok = s.s[key]
	log.Info("Get key:", key, value)
	return value, ok
}

// Set stores a new URL mapping in the file storage.
// Returns an error if the key is empty, value is empty, or if the key already exists.
func (s *FileStore) Set(ctx context.Context, key string, value string) error {
	log := logger.WithContext(ctx)

	if key == "" || value == "" {
		return fmt.Errorf("invalid input")
	}

	if _, ok := s.s[key]; ok {
		return appErrors.NewHTTPError(409, "conflict")
	}

	log.Info("Set key:", key, value)
	s.s[key] = value

	WriteToFile(key, value, s.pathToFile)

	return nil
}

// List returns all URL mappings currently stored in the file.
func (s *FileStore) List(ctx context.Context) (map[string]string, error) {
	return s.s, nil
}

// Ping always returns nil for file storage as it doesn't require connection checking.
func (s *FileStore) Ping(ctx context.Context) error {
	return nil
}

// SetMultiple stores multiple URL mappings in a batch operation.
func (s *FileStore) SetMultiple(ctx context.Context, items map[string]string) error {
	log := logger.WithContext(ctx)

	for k, v := range items {
		log.Info("SetMultiple", "key", k, "value", v)
		s.s[k] = v

		WriteToFile(k, v, s.pathToFile)
	}
	return nil
}

// Delete marks one or more URLs as deleted by setting their values to empty string.
func (s *FileStore) Delete(ctx context.Context, hashes []string) error {
	log := logger.WithContext(ctx)

	for _, v := range hashes {
		log.Info("Delete", "hash", v)
		s.s[v] = ""
	}
	return nil
}

// getDirFromPath extracts the directory path from a full file path.
func getDirFromPath(path string) (dir string) {
	parts := strings.Split(path, `/`)
	return strings.Join(parts[:len(parts)-1], "/")
}

// createDir creates the directory structure for the storage file if it doesn't exist.
func createDir(storagePath string) {
	log := logger.Get()

	d := getDirFromPath(storagePath)

	if err := os.MkdirAll(d, os.ModePerm); err != nil {
		log.Error("createDir", "storagePath", storagePath, "err", err)
		panic(err)
	}
}

// createFile creates a new storage file if it doesn't exist.
func createFile(storagePath string) {
	log := logger.Get()

	if _, err := os.Stat(storagePath); err == nil {
		// file exists
		return
	}

	data := []byte("")
	f, err := os.Create(storagePath)

	if err != nil {
		log.Error("createFile", "storagePath", storagePath, "err", err)
		panic(err)
	}
	defer f.Close()
	f.Write(data)
}

// itemURL represents the JSON structure used for file storage.
type itemURL struct {
	Hash string `json:"Hash"`
	URL  string `json:"OriginalURL"`
}

// lineToJSON parses a line from the storage file into an itemURL struct.
func lineToJSON(line string) itemURL {
	req := itemURL{}
	if err := json.Unmarshal([]byte(line), &req); err != nil {
		panic(err)
	}
	return req
}

// getFileContent reads and parses all URL mappings from the storage file.
func getFileContent(storagePath string) map[string]string {
	log := logger.Get()

	file, err := os.Open(storagePath)
	if err != nil {
		log.Error("getFileContent", "err", err)
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
		log.Error("getFileContent", "scanner.err", err)
	}

	return data
}

// NewFileStore creates and initializes a new file-based storage instance.
// It ensures the storage directory and file exist, and loads any existing data.
func NewFileStore(storagePath string) Storage {
	createDir(storagePath)
	createFile(storagePath)
	data := getFileContent(storagePath)

	return &FileStore{
		pathToFile: storagePath,
		s:          data,
	}
}

// WriteToFile appends a new URL mapping to the storage file in JSON format.
func WriteToFile(key string, value string, pathToFile string) error {
	log := logger.Get()

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

	log.Debug("WriteToFile", "json", json)
	return nil
}
