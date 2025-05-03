package randx

import (
	crand "crypto/rand"
	"io"
	"math/rand/v2"
)

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
	bytes := make([]byte, len)
	_, err := io.ReadFull(crand.Reader, bytes)
	if err != nil {
		panic(err)
	}
	return bytes
}
