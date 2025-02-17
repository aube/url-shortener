package store

import (
	"github.com/aube/url-shortener/internal/app/hasher"
)

var urlsMap = make(map[string]string)

func GetURLHash(id string) string {
	return urlsMap[id]
}

func SetURLHash(body []byte) string {
	hash := hasher.CalcHash(body)
	urlsMap[hash] = string(body)
	return hash
}
