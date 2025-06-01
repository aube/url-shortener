package hasher

import (
	"crypto/sha1"
	"fmt"

	"github.com/aube/url-shortener/internal/logger"
)

// CalcHash calculates a SHA-1 hash of the input data and returns a substring
func CalcHash(body []byte) string {
	hasher := sha1.New()
	log := logger.Get()

	_, err := hasher.Write(body)
	if err != nil {
		log.Error("CalcHash", "err", err)
		return ""
	}

	hashBytes := hasher.Sum(nil)
	hashString := fmt.Sprintf("%x", hashBytes)
	substringLength := min(10, len(hashString)-1)

	return hashString[:substringLength]
}
