package main

import (
	"bytes"
	"cryptopals/internal/cipherx"
	"cryptopals/internal/crack"
	"cryptopals/internal/randx"
	"encoding/base64"
	"fmt"
	"log"
	"os"
)

type Ch12Oracle struct {
	key       []byte
	plaintext []byte
}

func makeOracle(plaintext []byte, keySize int) Ch12Oracle {
	key := randx.RandBytes(keySize)
	return Ch12Oracle{key: key, plaintext: plaintext}
}

func (o Ch12Oracle) Encrypt(userPlaintext []byte) []byte {
	const bs = 16

	buf := make([]byte, 0, len(userPlaintext)+len(o.plaintext)+bs)
	buf = append(buf, userPlaintext...)
	buf = append(buf, o.plaintext...)
	buf = cipherx.Pcks7Padding(buf, bs)

	output, err := cipherx.EncryptAesEcb(buf, o.key)
	if err != nil {
		panic(err)
	}

	return output
}

func (o Ch12Oracle) CheckAnswer(answer []byte) bool {
	return bytes.Equal(answer, o.plaintext)
}

func crackBlockSize(o Ch12Oracle) int {
	buf := make([]byte, 256)
	var prevFirstBlock []byte
	for i := range buf {
		buf[i] = 'a'
		output := o.Encrypt(buf)
		if prevFirstBlock != nil && bytes.Equal(prevFirstBlock, output[:i]) {
			return i
		}
		prevFirstBlock = output[:i+1]
	}
	return 0
}

func detectEbc(o Ch12Oracle, bs int) bool {
	ciphertext := o.Encrypt([]byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"))
	return crack.HasRepeatingBlock(ciphertext, bs)
}

func detectLength(o Ch12Oracle, bs int) int {
	pad16Plaintext := make([]byte, bs)
	for i := range pad16Plaintext {
		pad16Plaintext[i] = 0x10
	}
	pad16Ciphertext := o.Encrypt(pad16Plaintext)[:bs]

	for initialPadLen := 0; initialPadLen < bs; initialPadLen++ {
		ct := o.Encrypt(make([]byte, initialPadLen))
		if len(ct)%bs != 0 {
			panic("length of ciphertext is not divisible by block size")
		}
		if bytes.Equal(ct[len(ct)-bs:], pad16Ciphertext) {
			return len(ct) - bs - initialPadLen
		}
	}

	panic("could not compute length")
}

func crackEbc(o Ch12Oracle, bs int, ctLen int) []byte {
	padLen := bs - 1
	paddedGuess := make([]byte, padLen+ctLen)

	for i := 0; i < ctLen; i++ {
		paddedGuessIndex := padLen + i

		blockChars := make(map[string]byte)
		for c := 0; c < 256; c++ {
			c := byte(c)
			paddedGuess[paddedGuessIndex] = c
			testBlock := paddedGuess[paddedGuessIndex+1-bs : paddedGuessIndex+1]
			ct := o.Encrypt(testBlock)
			blockChars[string(ct[:bs])] = c
		}

		shiftLen := (bs - 1) - (i % 16)
		ct := o.Encrypt(make([]byte, shiftLen))
		blockNum := i / bs
		c, exists := blockChars[string(ct[bs*blockNum:bs*(blockNum+1)])]
		if !exists {
			log.Fatalf("i=%v, blockNum=%v, shiftLen=%v, block does not exist in map\n", i, blockNum, shiftLen)
		}
		paddedGuess[paddedGuessIndex] = c
	}

	return paddedGuess[padLen:]
}

func main() {
	plaintext, err := base64.StdEncoding.DecodeString("Um9sbGluJyBpbiBteSA1LjAKV2l0aCBteSByYWctdG9wIGRvd24gc28gbXkgaGFpciBjYW4gYmxvdwpUaGUgZ2lybGllcyBvbiBzdGFuZGJ5IHdhdmluZyBqdXN0IHRvIHNheSBoaQpEaWQgeW91IHN0b3A/IE5vLCBJIGp1c3QgZHJvdmUgYnkK")
	if err != nil {
		panic(err)
	}

	oracle := makeOracle(plaintext, 16)

	bs := crackBlockSize(oracle)
	fmt.Printf("Block size: %v\n", bs)

	isEbc := detectEbc(oracle, bs)
	if isEbc {
		fmt.Printf("EBC detected, proceeding\n")
	} else {
		fmt.Printf("EBC not detected, exiting\n")
		os.Exit(1)
	}

	ctLen := detectLength(oracle, bs)
	fmt.Printf("Detected ciphertext length: %v\n", ctLen)

	answer := crackEbc(oracle, bs, ctLen)
	fmt.Printf("%q\n", answer)
	if oracle.CheckAnswer(answer) {
		fmt.Printf("Answer is correct\n")
	} else {
		fmt.Printf("Answer is incorrect\n")
	}
}
