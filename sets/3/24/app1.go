package main

import "cryptopals/internal/randx"

func encryptWithMersenneTwisterCipher(pt []byte) (uint16, []byte) {
	seed := RandUint16()
	cipher := MakeMersenneTwisterCipher(seed)
	ct := cipher.Encrypt(append(randx.RandBytes(randx.RandRange(5, 20)), pt...))
	return seed, ct
}
