package hasher

import (
	"crypto/sha1"
	"fmt"
)

func CalcHash(body []byte) string {
	// TODO:
	// var hash uint32 = crc32.ChecksumIEEE([]byte(data))

	hasher := sha1.New()
	hasher.Write(body)
	hashBytes := hasher.Sum(nil)
	hashString := fmt.Sprintf("%x", hashBytes)
	substringLength := min(10, len(hashString)-1)

	return hashString[:substringLength]
}
