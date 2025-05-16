package main

import (
	"cryptopals/internal/cipherx"
	"cryptopals/internal/crack"
	"cryptopals/internal/randx"
	"cryptopals/internal/util"
	"fmt"
)

func getCiphertexts() [][]byte {
	strings, err := util.ReadBase64ListFile("20.txt")
	if err != nil {
		panic(err)
	}
	cipher, err := cipherx.MakeAesCtr(randx.RandBytes(16))
	if err != nil {
		panic(err)
	}
	var ciphertexts [][]byte
	for _, s := range strings {
		ct := cipher.Encrypt(s, 0)
		ciphertexts = append(ciphertexts, ct)
	}
	return ciphertexts
}

func main() {
	ciphertexts := getCiphertexts()

	var longestLen int
	for i, ct := range ciphertexts {
		if i == 0 || len(ct) > longestLen {
			longestLen = len(ct)
		}
	}

	key := make([]byte, longestLen)
	for offset := range longestLen {
		buf := make([]byte, 0, len(ciphertexts))
		for _, ct := range ciphertexts {
			if offset < len(ct) {
				buf = append(buf, ct[offset])
			}
		}
		var result crack.CrackSingleXorResult
		if offset == 0 {
			result = crack.CrackSingleXorFirstCharacter(buf)
		} else {
			result = crack.CrackSingleXor(buf)
		}
		key[offset] = result.Key
	}

	for _, ct := range ciphertexts {
		pt := make([]byte, len(ct))
		cipherx.XorBytes(pt, ct, key[:len(ct)])
		fmt.Printf("%q\n", pt)
	}
}
