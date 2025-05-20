package hasher

import (
	"fmt"
)

func Example() {
	fakeAddress := "http://test.test/test"

	hash := CalcHash([]byte(fakeAddress))
	conflictHash := CalcHash([]byte("conflict"))

	fmt.Println("hash == conflictHash", hash == conflictHash)
}
