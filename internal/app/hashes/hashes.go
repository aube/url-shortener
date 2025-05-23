package hashes

import (
	"crypto/sha1"
	"fmt"
)

var urlsMap = make(map[string]string)

func hashCalc(body []byte) string {
	// TODO:
	// var hash uint32 = crc32.ChecksumIEEE([]byte(data))

	hasher := sha1.New()
	hasher.Write(body)
	hashBytes := hasher.Sum(nil)
	hashString := fmt.Sprintf("%x", hashBytes)
	substringLength := min(10, len(hashString)-1)

	return hashString[:substringLength]
}

func GetURLHash(id string) string {
	return urlsMap[id]
}

func SetURLHash(body []byte) string {
	hash := hashCalc(body)
	urlsMap[hash] = string(body)
	return hash
}
