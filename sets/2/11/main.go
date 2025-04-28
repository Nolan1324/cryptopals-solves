package main

import (
	"cryptopals/internal/cipherx"
	"cryptopals/internal/crack"
	"fmt"
	"math/rand/v2"
)

func randRange(min, max int) int {
	return rand.IntN(max-min) + min
}

func randByte() byte {
	return byte(rand.UintN(256))
}

func randBytes(len int) []byte {
	var bytes []byte
	for i := 0; i < len; i++ {
		bytes = append(bytes, randByte())
	}
	return bytes
}

func encryptionOracle(plaintext []byte) ([]byte, bool) {
	const bs = 16

	appendRandBytes := func(buf []byte) {
		for i := 0; i < randRange(5, 11); i++ {
			buf = append(buf, randByte())
		}
	}

	var buf []byte
	appendRandBytes(buf)
	buf = append(buf, plaintext...)
	appendRandBytes(buf)
	buf = cipherx.Pcks7Padding(buf, bs)

	key := randBytes(bs)

	isEcb := rand.IntN(2) == 0
	var output []byte
	var err error
	if isEcb {
		output, err = cipherx.EncryptAesEcb(buf, key)
	} else {
		output, err = cipherx.EncryptAesCbc(buf, key, randBytes(bs))
	}

	if err != nil {
		panic(err)
	}

	return output, isEcb
}

func main() {
	ciphertext, isEcb := encryptionOracle([]byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"))
	if crack.HasRepeatingBlock(ciphertext, 16) {
		fmt.Printf("Detected ECB\n")
	} else {
		fmt.Printf("Detected CBC\n")
	}
	if isEcb {
		fmt.Printf("True answer: ECB\n")
	} else {
		fmt.Printf("True answer: CBC\n")
	}
}
