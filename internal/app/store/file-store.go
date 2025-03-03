package store

import (
	"bufio"
	"encoding/json"
	"os"
	"strings"

	"github.com/aube/url-shortener/internal/logger"
)

var storagePathFile string

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

func putFileIntoMem(storagePath string) {
	file, err := os.Open(storagePath)
	if err != nil {
		logger.Println(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			json := lineToJSON(line)
			SetValue(json.Hash, json.URL)
		}
	}

	if err = scanner.Err(); err != nil {
		logger.Println(err)
	}
}

func NewFileStore(storagePath string) {
	createDir(storagePath)
	createFile(storagePath)
	putFileIntoMem(storagePath)
	storagePathFile = storagePath
}

func WriteToFile(key string, value string) error {
	f, err := os.OpenFile(storagePathFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
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
