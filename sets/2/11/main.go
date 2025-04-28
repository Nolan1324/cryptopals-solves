package main

import (
	"cryptopals/internal/cipherx"
	"cryptopals/internal/crack"
	"cryptopals/internal/randx"
	"fmt"
	"math/rand/v2"
)

type Ch11Oracle struct {
	keySize int
	isEcb   bool
}

func makeOracle(keySize int) Ch11Oracle {
	return Ch11Oracle{keySize: keySize, isEcb: rand.IntN(2) == 0}
}

func (o Ch11Oracle) encrypt(plaintext []byte) []byte {
	const bs = 16

	appendRandBytes := func(buf []byte) {
		for i := 0; i < randx.RandRange(5, 11); i++ {
			buf = append(buf, randx.RandByte())
		}
	}

	var buf []byte
	appendRandBytes(buf)
	buf = append(buf, plaintext...)
	appendRandBytes(buf)
	buf = cipherx.Pcks7Padding(buf, bs)

	key := randx.RandBytes(o.keySize)

	var output []byte
	var err error
	if o.isEcb {
		output, err = cipherx.EncryptAesEcb(buf, key)
	} else {
		output, err = cipherx.EncryptAesCbc(buf, key, randx.RandBytes(bs))
	}

	if err != nil {
		panic(err)
	}

	return output
}

func main() {
	oracle := makeOracle(16)
	ciphertext := oracle.encrypt([]byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"))
	if crack.HasRepeatingBlock(ciphertext, 16) {
		fmt.Printf("Detected ECB\n")
	} else {
		fmt.Printf("Detected CBC\n")
	}
	if oracle.isEcb {
		fmt.Printf("True answer: ECB\n")
	} else {
		fmt.Printf("True answer: CBC\n")
	}
}
