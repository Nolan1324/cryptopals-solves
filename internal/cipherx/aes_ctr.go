package cipherx

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
)

type count = uint64

const sizeOfCount = 8

const bs = 16

type aesCtr struct {
	cipher cipher.Block
}

type AesCtr interface {
	Encrypt(plaintext []byte, nonce count) []byte
	Decrypt(ciphertext []byte, nonce count) []byte
}

func NewAesCtr(key []byte) (AesCtr, error) {
	cipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return aesCtr{cipher: cipher}, nil
}

func (a aesCtr) Encrypt(pt []byte, nonce count) []byte {
	ct := make([]byte, len(pt))
	keyBlock := make([]byte, bs)
	count := nonce
	// Encrypt each full block
	for blockIdx := range len(pt) / bs {
		i := blockIdx * bs
		a.getKeyBlock(keyBlock, count)
		XorBytes(ct[i:i+bs], pt[i:i+bs], keyBlock)
		count++
	}
	// Encrypt the remaining text that is less than the block size, if any
	if len(pt)%bs != 0 {
		i := (len(pt) / bs) * bs
		a.getKeyBlock(keyBlock, count)
		XorBytes(ct[i:], pt[i:], keyBlock[:len(ct[i:])])
	}
	return ct
}

func (a aesCtr) Decrypt(ct []byte, nonce count) []byte {
	// CTR simply XORs the keystream with the text, so decrypt is the same as encrypt
	return a.Encrypt(ct, nonce)
}

func (a aesCtr) getKeyBlock(dst []byte, count count) {
	block := make([]byte, bs)
	binary.LittleEndian.PutUint64(block[bs-sizeOfCount:], count)
	a.cipher.Encrypt(dst, block)
}
