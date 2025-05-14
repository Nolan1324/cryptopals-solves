package main

import (
	"cryptopals/internal/crack"
	"encoding/base64"
	"fmt"
	"os"
)

func main() {
	plaintext, err := base64.StdEncoding.DecodeString("Um9sbGluJyBpbiBteSA1LjAKV2l0aCBteSByYWctdG9wIGRvd24gc28gbXkgaGFpciBjYW4gYmxvdwpUaGUgZ2lybGllcyBvbiBzdGFuZGJ5IHdhdmluZyBqdXN0IHRvIHNheSBoaQpEaWQgeW91IHN0b3A/IE5vLCBJIGp1c3QgZHJvdmUgYnkK")
	if err != nil {
		panic(err)
	}

	oracle := makeOracle(plaintext, 16)

	bs := crack.DetectBlockSize(oracle.Encrypt)
	fmt.Printf("Detected block size: %v\n", bs)
	if bs != 16 {
		fmt.Printf("Incorrect block size detected\n")
		os.Exit(1)
	}

	isEcb := crack.DetectEcbMode(oracle.Encrypt, bs)
	if isEcb {
		fmt.Printf("Detected ECB mode\n")
	} else {
		fmt.Printf("Did not detect ECB mode\n")
		os.Exit(1)
	}

	ctLen := crack.DetectEcbLength(oracle.Encrypt, bs)
	fmt.Printf("Detected ciphertext length: %v\n", ctLen)

	answer := crack.CrackEcb(oracle.Encrypt, bs, ctLen)
	fmt.Printf("%q\n", answer)
	if oracle.CheckAnswer(answer) {
		fmt.Printf("Answer is correct\n")
	} else {
		fmt.Printf("Answer is incorrect\n")
	}
}
