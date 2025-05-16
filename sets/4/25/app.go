package main

import (
	"bytes"
	"cryptopals/internal/cipherx"
	"cryptopals/internal/randx"
	"cryptopals/internal/util"
)

type Application struct {
	ciphertext []byte
	ctr        cipherx.AesCtr
}

func MakeApplication() Application {
	ctEcb := util.ReadBase64File("25.txt")
	pt, err := cipherx.DecryptAesEcb(ctEcb, []byte("YELLOW SUBMARINE"))
	if err != nil {
		panic(err)
	}

	ctr, err := cipherx.MakeAesCtr(randx.RandBytes(16))
	if err != nil {
		panic(err)
	}
	return Application{
		ciphertext: ctr.Encrypt(pt, 0),
		ctr:        ctr,
	}
}

// EditByte changes the encrypted byte at offset to newByte
func (a Application) EditByte(offset int, newByte byte) {
	a.ciphertext[offset] = newByte ^ a.ctr.GetKeystreamByte(0, uint64(offset))
}

// ReadCiphertext reads the entire current ciphertext
func (a Application) ReadCiphertext() []byte {
	return bytes.Clone(a.ciphertext)
}
