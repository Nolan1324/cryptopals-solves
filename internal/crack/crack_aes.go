package crack

import (
	"bytes"
	"log"
)

func HasRepeatingBlock(bytes []byte, bs int) bool {
	block_map := make(map[string]int)
	for i := 0; i < len(bytes); i += bs {
		block := string(bytes[i : i+bs])
		block_map[block]++
	}
	for _, v := range block_map {
		if v > 1 {
			return true
		}
	}
	return false
}

type EncryptFunc func([]byte) []byte

// DetectBlockSize detects the block size of encryptFunc.
// encryptFunc is expected to compute encryptFunc(plaintext) = BLOCK_CIPHER(prefix || plaintext || suffix)
func DetectBlockSize(encryptFunc EncryptFunc) int {
	originalLen := len(encryptFunc(nil))
	for i := 1; i <= 16; i++ {
		newLen := len(encryptFunc(make([]byte, i)))
		if newLen > originalLen {
			return newLen - originalLen
		}
	}
	return 0
}

// DetectEcbMode detects if a block cipher with block size bs is using ECB mode.
// encryptFunc is expected to compute encryptFunc(plaintext) = BLOCK_CIPHER(prefix || plaintext || suffix)
func DetectEcbMode(encryptFunc EncryptFunc, bs int) bool {
	return HasRepeatingBlock(encryptFunc(make([]byte, bs*4)), bs)
}

// DetectEcbBlockSizeOneShot detects if encryptFunc is using ECB and returns the block size if so.
// Returns 0 if ECB is not detected. This function only makes one call to encryptFunc.
// encryptFunc is expected to compute encryptFunc(plaintext) = AES_ECB(prefix || plaintext || suffix)
// Tries all block sizes in [min, max] and returns the first one that results in a pair of consecutive identical blocks in the ciphertext.
// min should typically be at least 8.
func DetectEcbBlockSizeOneShot(encryptFunc EncryptFunc, min int, max int) int {
	buf := make([]byte, max*3)
	output := encryptFunc(buf)
	for bs := min; bs <= max; bs++ {
		if len(output)%bs != 0 {
			continue
		}
		// Find a pair of consecutive equal blocks
		for i := bs; i < len(output); i += bs {
			if bytes.Equal(output[i-bs:i], output[i:i+bs]) {
				return bs
			}
		}
	}
	return 0
}

// DetectEcbLength detects the target text length, given a function encryptFunc of the form
// encryptFunc(plaintext) = AES_ECB(plaintext || target_text) with a known ECB block size of bs
func DetectEcbLength(encryptFunc EncryptFunc, bs int) int {
	pad16Plaintext := make([]byte, bs)
	for i := range pad16Plaintext {
		pad16Plaintext[i] = 0x10
	}
	pad16Ciphertext := encryptFunc(pad16Plaintext)[:bs]

	for initialPadLen := range bs {
		ct := encryptFunc(make([]byte, initialPadLen))
		if len(ct)%bs != 0 {
			panic("length of ciphertext is not divisible by block size")
		}
		if bytes.Equal(ct[len(ct)-bs:], pad16Ciphertext) {
			return len(ct) - bs - initialPadLen
		}
	}

	panic("could not compute length")
}

// CrackEcb cracks the target text, given a function encryptFunc of the form
// encryptFunc(plaintext) = AES_ECB(plaintext || target_text) with a known ECB block size of bs
// and target text length of ctLen
func CrackEcb(encryptFunc EncryptFunc, bs int, ctLen int) []byte {
	padLen := bs - 1
	paddedGuess := make([]byte, padLen+ctLen)

	for i := range ctLen {
		paddedGuessIndex := padLen + i

		blockChars := make(map[string]byte)
		for c := range 256 {
			c := byte(c)
			paddedGuess[paddedGuessIndex] = c
			testBlock := paddedGuess[paddedGuessIndex+1-bs : paddedGuessIndex+1]
			ct := encryptFunc(testBlock)
			blockChars[string(ct[:bs])] = c
		}

		shiftLen := (bs - 1) - (i % 16)
		ct := encryptFunc(make([]byte, shiftLen))
		blockNum := i / bs
		c, exists := blockChars[string(ct[bs*blockNum:bs*(blockNum+1)])]
		if !exists {
			log.Fatalf("i=%v, blockNum=%v, shiftLen=%v, block does not exist in map\n", i, blockNum, shiftLen)
		}
		paddedGuess[paddedGuessIndex] = c
	}

	return paddedGuess[padLen:]
}
