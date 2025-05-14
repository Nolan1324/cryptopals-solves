package main

import (
	"bytes"
	"cryptopals/internal/crack"
	"encoding/base64"
	"fmt"
	"os"
)

// findZeroBlockCiphertext finds the ciphertext associated with a block of all zeroes.
func findZeroBlockCiphertext(oracle Ch14Oracle, bs int) []byte {
	output := oracle.Encrypt(make([]byte, bs*4))
	for i := bs; i < len(output); i += bs {
		if bytes.Equal(output[i-bs:i], output[i:i+bs]) {
			return output[i-bs : i]
		}
	}
	panic("Zero block ciphertext not found")
}

// detectPrefixLength detects the length of the prefix within the encryption oracle for challenge 14.
func detectPrefixLength(oracle Ch14Oracle, bs int) int {
	zeroBlockCt := findZeroBlockCiphertext(oracle, bs)

	// Keep encrypting with more zero bytes inserted after the prefix until we create the zero block.
	for numZeroBytes := bs; numZeroBytes < 2*bs; numZeroBytes++ {
		// We append a 'x' at the end of the zero bytes to seperate them from the target text,
		// in case the target text happens to start with \x00
		output := oracle.Encrypt(append(make([]byte, numZeroBytes), 'x'))
		// Check if the zero block is present
		for blockOffset := 0; blockOffset < len(output); blockOffset += bs {
			if bytes.Equal(zeroBlockCt, output[blockOffset:blockOffset+bs]) {
				// Other than the zero bytes that formed the zero block, numZeroBytes - bs
				// joined with the end of the prefix to complete the final prefix block
				// So the prefix length is blockOffset (where the prefix block stops and the zero block starts)
				// minus these extra zero bytes in the prefix block
				return blockOffset - (numZeroBytes - bs)
			}
		}
	}

	panic("Prefix length not detected")
}

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

	prefixLength := detectPrefixLength(oracle, bs)
	fmt.Printf("Detected prefix length %v\n", prefixLength)
	if prefixLength != len(oracle.prefix) {
		fmt.Printf("Correct prefix length is %v\n, which does not match\n", len(oracle.prefix))
		os.Exit(1)
	}

	encrypt := func(plaintext []byte) []byte {
		startBlockIndex := prefixLength/bs + 1
		prefixPaddingLen := bs - (prefixLength % bs)
		return oracle.Encrypt(append(make([]byte, prefixPaddingLen), plaintext...))[startBlockIndex*bs:]
	}

	ctLen := crack.DetectEcbLength(encrypt, bs)
	fmt.Printf("Detected ciphertext length: %v\n", ctLen)

	answer := crack.CrackEcb(encrypt, bs, ctLen)
	fmt.Printf("%q\n", answer)
	if oracle.CheckAnswer(answer) {
		fmt.Printf("Answer is correct\n")
	} else {
		fmt.Printf("Answer is incorrect\n")
	}
}
