package randx

import "math/rand/v2"

func RandRange(min, max int) int {
	return rand.IntN(max-min) + min
}

func RandByte() byte {
	return byte(rand.UintN(256))
}

func RandBytes(len int) []byte {
	var bytes []byte
	for i := 0; i < len; i++ {
		bytes = append(bytes, RandByte())
	}
	return bytes
}
