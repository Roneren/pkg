// +build gofuzz

package phpsession

func Fuzz(data []byte) int {
	decoder := NewPhpDecoder(string(data))
	_, err := decoder.Decode()

	if err != nil {
		return 0
	}

	return 1
}
