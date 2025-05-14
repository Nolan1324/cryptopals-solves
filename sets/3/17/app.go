package main

import (
	"cryptopals/internal/cipherx"
	"cryptopals/internal/randx"
	"encoding/base64"
	"math/rand/v2"
)

var strings = [...]string{
	"MDAwMDAwTm93IHRoYXQgdGhlIHBhcnR5IGlzIGp1bXBpbmc=",
	"MDAwMDAxV2l0aCB0aGUgYmFzcyBraWNrZWQgaW4gYW5kIHRoZSBWZWdhJ3MgYXJlIHB1bXBpbic=",
	"MDAwMDAyUXVpY2sgdG8gdGhlIHBvaW50LCB0byB0aGUgcG9pbnQsIG5vIGZha2luZw==",
	"MDAwMDAzQ29va2luZyBNQydzIGxpa2UgYSBwb3VuZCBvZiBiYWNvbg==",
	"MDAwMDA0QnVybmluZyAnZW0sIGlmIHlvdSBhaW4ndCBxdWljayBhbmQgbmltYmxl",
	"MDAwMDA1SSBnbyBjcmF6eSB3aGVuIEkgaGVhciBhIGN5bWJhbA==",
	"MDAwMDA2QW5kIGEgaGlnaCBoYXQgd2l0aCBhIHNvdXBlZCB1cCB0ZW1wbw==",
	"MDAwMDA3SSdtIG9uIGEgcm9sbCwgaXQncyB0aW1lIHRvIGdvIHNvbG8=",
	"MDAwMDA4b2xsaW4nIGluIG15IGZpdmUgcG9pbnQgb2g=",
	"MDAwMDA5aXRoIG15IHJhZy10b3AgZG93biBzbyBteSBoYWlyIGNhbiBibG93",
}

type Application struct {
	key []byte
}

func makeApplication() Application {
	return Application{key: randx.RandBytes(16)}
}

func (a Application) CreateCiphertext() ([]byte, []byte) {
	iv := randx.RandBytes(16)
	pt, err := base64.StdEncoding.DecodeString(strings[rand.IntN(len(strings))])
	if err != nil {
		panic(err)
	}
	pt = cipherx.AddPkcs7Padding([]byte(pt), 16)
	ct, err := cipherx.EncryptAesCbc(pt, a.key, iv)
	if err != nil {
		panic(err)
	}
	return ct, iv
}

func (a Application) HasValidPadding(ciphertext []byte, iv []byte) bool {
	pt, err := cipherx.DecryptAesCbc(ciphertext, a.key, iv)
	if err != nil {
		return false
	}
	_, err = cipherx.RemovePkcs7Padding(pt)
	// if err == nil {
	// 	log.Printf("pt=%q", pt)
	// }
	return err == nil
}
