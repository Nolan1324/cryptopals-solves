package crack

import (
	"cryptopals/internal/cipherx"
	"cryptopals/internal/histogram"
	"log"
	"sort"

	"gonum.org/v1/gonum/mat"
)

var frequencyVec *mat.VecDense

func initFrequencyVec() {
	if frequencyVec != nil {
		return
	}
	var frequencies = [26]float64{0.082, 0.015, 0.028, 0.043, 0.127, 0.022, 0.02, 0.061, 0.07, 0.0015, 0.0077, 0.04, 0.024, 0.067, 0.075, 0.019, 0.00095, 0.06, 0.063, 0.091, 0.028, 0.0098, 0.024, 0.0015, 0.02, 0.00074}
	frequencyVec = mat.NewVecDense(26, frequencies[:])
}

func CrackSingleXor(buf []byte) ([]byte, byte, float64) {
	var bestGuess []byte
	var bestScore float64
	var bestKey byte
	for i := 0; i < 256; i++ {
		guess := make([]byte, len(buf))
		cipherx.XorByte(guess, buf, byte(i))

		score := histogram.Score(guess)
		if score > bestScore {
			bestScore = score
			bestGuess = guess
			bestKey = byte(i)
		}
	}
	return bestGuess, bestKey, bestScore
}

type KeySize struct {
	Size         int
	EditDistance float32
}

func GuessXorKeySizes(buf []byte, min int, max int) []KeySize {
	var keySizes []KeySize
	for keySize := min; keySize <= max; keySize++ {
		dist, err := cipherx.EditDistance(buf[0:keySize], buf[keySize:2*keySize])
		if err != nil {
			log.Fatalln(err)
		}
		distNorm := float32(dist) / float32(keySize)
		keySizes = append(keySizes, KeySize{Size: keySize, EditDistance: distNorm})
	}

	sort.Slice(keySizes, func(i, j int) bool {
		return keySizes[i].EditDistance < keySizes[j].EditDistance
	})

	return keySizes
}

func CrackRepeatingKeyXorGivenKeySize(buf []byte, keySize int) []byte {
	key := make([]byte, 0, keySize)
	for offset := 0; offset < keySize; offset++ {
		newBuf := make([]byte, 0, len(buf)/keySize+1)
		for i := offset; i < len(buf); i += keySize {
			newBuf = append(newBuf, buf[i])
		}
		_, keyByte, _ := CrackSingleXor(newBuf)
		key = append(key, keyByte)
	}
	return key
}

func CrackRepeatingKeyXor(buf []byte, minKeySize int, maxKeySize int, topNumKeys int) []byte {
	keySizes := GuessXorKeySizes(buf, minKeySize, maxKeySize)
	keySizes = keySizes[:topNumKeys]

	var bestGuess []byte
	var bestScore float64
	for _, keySize := range keySizes {
		key := CrackRepeatingKeyXorGivenKeySize(buf, keySize.Size)
		guess := make([]byte, len(buf))
		cipherx.RepeatingKeyXor(guess, buf, key)
		score := histogram.Score(guess)
		if score > bestScore {
			bestGuess = guess
			bestScore = score
		}
	}

	return bestGuess
}
