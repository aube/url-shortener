package hasher

import (
	"crypto/sha1"
	"fmt"
)

// CalcHash calculates a SHA-1 hash of the input data and returns a substring
func CalcHash(body []byte) string {
	hasher := sha1.New()
	hasher.Write(body)
	hashBytes := hasher.Sum(nil)
	hashString := fmt.Sprintf("%x", hashBytes)
	substringLength := min(10, len(hashString)-1)

	return hashString[:substringLength]
}
