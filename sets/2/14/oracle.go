package main

import (
	"bytes"
	"cryptopals/internal/cipherx"
	"cryptopals/internal/randx"
)

type Ch14Oracle struct {
	key       []byte
	prefix    []byte
	plaintext []byte
}

func makeOracle(plaintext []byte, keySize int) Ch14Oracle {
	key := randx.RandBytes(keySize)
	prefix := randx.RandBytes(randx.RandRange(4, 128))
	return Ch14Oracle{key: key, prefix: prefix, plaintext: plaintext}
}

func (o Ch14Oracle) Encrypt(userPlaintext []byte) []byte {
	const bs = 16

	buf := make([]byte, 0, len(o.prefix)+len(userPlaintext)+len(o.plaintext)+bs)
	buf = append(buf, o.prefix...)
	buf = append(buf, userPlaintext...)
	buf = append(buf, o.plaintext...)
	buf = cipherx.AddPcks7Padding(buf, bs)

	output, err := cipherx.EncryptAesEcb(buf, o.key)
	if err != nil {
		panic(err)
	}

	return output
}

func (o Ch14Oracle) CheckAnswer(answer []byte) bool {
	return bytes.Equal(answer, o.plaintext)
}
