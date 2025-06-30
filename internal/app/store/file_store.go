package store

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	appErrors "github.com/aube/url-shortener/internal/app/apperrors"
	"github.com/aube/url-shortener/internal/app/ctxkeys"
	"github.com/aube/url-shortener/internal/logger"
)

// FileStore is a file-based implementation of the Storage interface
// that persists URL mappings to disk in JSON format.
type FileStore struct {
	urls      map[string]string
	users     map[string]string
	usersFile string
	urlsFile  string
}

// itemURL represents the JSON structure used for file storage.
type itemURL struct {
	Hash string `json:"Hash"`
	URL  string `json:"OriginalURL"`
}

// itemUser represents the JSON structure used for file storage.
type itemUser struct {
	Hash   string `json:"Hash"`
	UserID string `json:"UserID"`
}

// Get retrieves a URL by its shortened key from the file storage.
// Returns the URL and true if found, empty string and false otherwise.
func (s *FileStore) Get(ctx context.Context, key string) (value string, ok bool) {
	log := logger.WithContext(ctx)

	value, ok = s.urls[key]
	if ok {
		log.Info("Get", "key", key, "value", value)
		return value, ok
	}

	return "", false
}

// Get retrieves a URL by its shortened key from the file storage.
// Returns the URL and true if found, empty string and false otherwise.
func (s *FileStore) GetByUser(ctx context.Context, key string) (value string, ok bool) {
	log := logger.WithContext(ctx)
	userID := ctx.Value(ctxkeys.UserIDKey).(string)

	value, ok = s.urls[key]
	if ok && s.users[key] == userID {
		log.Info("Get", "key", key, "value", value)
		return value, ok
	}

	return "", false
}

// Set stores a new URL mapping in the file storage.
// Returns an error if the key is empty, value is empty, or if the key already exists.
func (s *FileStore) Set(ctx context.Context, key string, value string) error {
	log := logger.WithContext(ctx)
	userID := ctx.Value(ctxkeys.UserIDKey).(string)

	if key == "" || value == "" {
		return fmt.Errorf("invalid input")
	}

	if _, ok := s.urls[key]; ok {
		return appErrors.NewHTTPError(409, "conflict")
	}

	log.Info("Set key:", key, value, "UserID", userID)
	s.urls[key] = value
	s.users[key] = userID

	json0, err := json.Marshal(itemURL{Hash: key, URL: value})
	if err != nil {
		return err
	}
	err = WriteToFile(s.urlsFile, json0)
	if err != nil {
		return err
	}

	json1, err := json.Marshal(itemUser{Hash: key, UserID: userID})
	if err != nil {
		return err
	}
	err = WriteToFile(s.usersFile, json1)
	if err != nil {
		return err
	}

	return nil
}

// List returns all URL mappings currently stored in the file.
func (s *FileStore) List(ctx context.Context) (map[string]string, error) {
	userID := ctx.Value(ctxkeys.UserIDKey).(string)
	urls := make(map[string]string)

	for hash, url := range s.urls {
		if s.users[hash] == userID {
			urls[hash] = url
		}
	}
	return urls, nil
}

// Ping always returns nil for file storage as it doesn't require connection checking.
func (s *FileStore) Ping(ctx context.Context) error {
	return nil
}

// SetMultiple stores multiple URL mappings in a batch operation.
func (s *FileStore) SetMultiple(ctx context.Context, items map[string]string) error {
	log := logger.WithContext(ctx)
	userID := ctx.Value(ctxkeys.UserIDKey).(string)

	for key, value := range items {
		log.Info("Set key:", key, value)
		s.urls[key] = value
		s.users[key] = userID

		json0, err := json.Marshal(itemURL{Hash: key, URL: value})
		if err != nil {
			return err
		}
		err = WriteToFile(s.urlsFile, json0)
		if err != nil {
			return err
		}

		json1, err := json.Marshal(itemUser{Hash: key, UserID: value})
		if err != nil {
			return err
		}
		err = WriteToFile(s.usersFile, json1)
		if err != nil {
			return err
		}
	}

	return nil
}

// Delete marks one or more URLs as deleted by setting their values to empty string.
func (s *FileStore) Delete(ctx context.Context, hashes []string) error {
	log := logger.WithContext(ctx)
	userID := ctx.Value(ctxkeys.UserIDKey).(string)

	for _, hash := range hashes {
		if s.users[hash] == userID {
			log.Info("Delete", "hash", hash)
			delete(s.urls, hash)
		}
	}
	return nil
}

// Stats select amount of urls and users from database.
func (s *FileStore) Stats(ctx context.Context) (int, int, error) {
	log := logger.WithContext(ctx)
	urls := len(s.urls)
	users := make(map[string]string)

	for _, user := range s.users {
		users[user] = ""
	}

	log.Info("Stats", "urls", urls, "users", len(users))
	return urls, len(users), nil
}

// createDir creates the directory structure for the storage file if it doesn't exist.
func createDir(storagePath string) {
	log := logger.Get()

	if err := os.MkdirAll(storagePath, os.ModePerm); err != nil {
		log.Error("createDir", "storagePath", storagePath, "err", err)
		panic(err)
	}
}

// createFile creates a new storage file if it doesn't exist.
func createFile(filePath string) error {
	log := logger.Get()

	if _, err := os.Stat(filePath); err == nil {
		// file exists
		return nil
	}

	f, err := os.Create(filePath)

	if err != nil {
		log.Error("createFile", "filePath", filePath, "err", err)
		panic(err)
	}
	defer f.Close()

	_, err = f.Write([]byte(""))
	if err != nil {
		log.Error("createFile2", "err", err)
		return err
	}
	return nil
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
	log := logger.Get()
	dir := filepath.Dir(storagePath)

	urlsFile := filepath.Join(storagePath)
	usersFile := filepath.Join(dir, "users_list.txt")

	createDir(dir)
	log.Info("NewFileStore", "createDir", storagePath)

	err := createFile(usersFile)

	if err != nil {
		log.Error("NewFileStore", "create usersFile", err)
		panic(fmt.Errorf("can't create file store: %w", err))
	}

	err = createFile(urlsFile)
	if err != nil {
		log.Error("NewFileStore", "create urlsFile", err)
		panic(fmt.Errorf("can't create file store: %w", err))
	}

	urls := getFileContent(urlsFile)
	users := getFileContent(usersFile)

	return &FileStore{
		usersFile: usersFile,
		urlsFile:  urlsFile,
		urls:      urls,
		users:     users,
	}
}

// WriteToFile appends a new URL mapping to the storage file in JSON format.
func WriteToFile(pathToFile string, json []byte) error {
	log := logger.Get()

	f, err := os.OpenFile(pathToFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err = f.WriteString(string(json) + "\n"); err != nil {
		return err
	}

	log.Debug("WriteToFile", "json", json)
	return nil
}
