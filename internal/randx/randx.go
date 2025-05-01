package randx

import "math/rand/v2"

// RandRange returns a random integer in the half-open interval [min, max). It panics if max < min.
func RandRange(min, max int) int {
	return rand.IntN(max-min) + min
}

// RandByte returns a random byte.
func RandByte() byte {
	return byte(rand.UintN(256))
}

// RandBytes returns len random bytes
func RandBytes(len int) []byte {
	var bytes []byte
	for i := 0; i < len; i++ {
		bytes = append(bytes, RandByte())
	}
	return bytes
}
