package cipherx

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
)

type Count = uint64

const sizeOfCount = 8

const bs = 16

type AesCtr struct {
	cipher cipher.Block
}

func MakeAesCtr(key []byte) (AesCtr, error) {
	cipher, err := aes.NewCipher(key)
	if err != nil {
		return AesCtr{}, err
	}
	return AesCtr{cipher: cipher}, nil
}

func (a AesCtr) Encrypt(pt []byte, nonce Count) []byte {
	ct := make([]byte, len(pt))
	keyBlock := make([]byte, bs)
	count := nonce
	// Encrypt each full block
	for blockIdx := range len(pt) / bs {
		i := blockIdx * bs
		a.GetKeystreamBlock(keyBlock, count)
		XorBytes(ct[i:i+bs], pt[i:i+bs], keyBlock)
		count++
	}
	// Encrypt the remaining text that is less than the block size, if any
	if len(pt)%bs != 0 {
		i := (len(pt) / bs) * bs
		a.GetKeystreamBlock(keyBlock, count)
		XorBytes(ct[i:], pt[i:], keyBlock[:len(ct[i:])])
	}
	return ct
}

func (a AesCtr) Decrypt(ct []byte, nonce Count) []byte {
	// CTR simply XORs the keystream with the text, so decrypt is the same as encrypt
	return a.Encrypt(ct, nonce)
}

func (a AesCtr) GetKeystreamBlock(dst []byte, count Count) {
	block := make([]byte, bs)
	binary.LittleEndian.PutUint64(block[bs-sizeOfCount:], count)
	a.cipher.Encrypt(dst, block)
}

func (a AesCtr) GetKeystreamByte(nonce Count, offset Count) byte {
	block := make([]byte, bs)
	a.GetKeystreamBlock(block, nonce+offset/bs)
	return block[offset%bs]
}
