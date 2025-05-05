package main

import (
	"cryptopals/internal/cipherx"
	"time"
)

func generatePasswordResetToken() []byte {
	rng := cipherx.NewMersenneTwister(uint32(time.Now().Unix()))
	token := make([]byte, 16)
	for i := range token {
		token[i] = byte(rng.Rand())
	}
	return token
}
