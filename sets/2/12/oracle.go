package main

import (
	"bytes"
	"cryptopals/internal/cipherx"
	"cryptopals/internal/randx"
)

type Ch12Oracle struct {
	key       []byte
	plaintext []byte
}

func makeOracle(plaintext []byte, keySize int) Ch12Oracle {
	key := randx.RandBytes(keySize)
	return Ch12Oracle{key: key, plaintext: plaintext}
}

func (o Ch12Oracle) Encrypt(userPlaintext []byte) []byte {
	const bs = 16

	buf := make([]byte, 0, len(userPlaintext)+len(o.plaintext)+bs)
	buf = append(buf, userPlaintext...)
	buf = append(buf, o.plaintext...)
	buf = cipherx.AddPcks7Padding(buf, bs)

	output, err := cipherx.EncryptAesEcb(buf, o.key)
	if err != nil {
		panic(err)
	}

	return output
}

func (o Ch12Oracle) CheckAnswer(answer []byte) bool {
	return bytes.Equal(answer, o.plaintext)
}
