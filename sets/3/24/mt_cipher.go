package main

import "cryptopals/internal/cipherx"

type MersenneTwisterCipher struct {
	key uint16
}

func MakeMersenneTwisterCipher(key uint16) MersenneTwisterCipher {
	return MersenneTwisterCipher{key: key}
}

func (c MersenneTwisterCipher) Encrypt(pt []byte) []byte {
	rng := cipherx.NewMersenneTwister(uint32(c.key))
	ct := make([]byte, len(pt))
	for i, b := range pt {
		ct[i] = b ^ byte(rng.Rand())
	}
	return ct
}

func (c MersenneTwisterCipher) Decrypt(ct []byte) []byte {
	return c.Encrypt(ct)
}
