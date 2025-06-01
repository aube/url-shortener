package hasher

func Example() {
	fakeAddress := "http://test.test/test"

	CalcHash([]byte(fakeAddress))
}
