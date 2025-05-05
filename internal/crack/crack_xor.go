package crack

import (
	"cryptopals/internal/cipherx"
	"cryptopals/internal/histogram"
	"sort"
)

type CrackSingleXorResult struct {
	Guess          []byte
	Score          float64
	Key            byte
	GuessHistogram histogram.Histogram
}

func CrackSingleXor(buf []byte) CrackSingleXorResult {
	return crackSingleXor(buf, nil, histogram.Score)
}

func CrackSingleXorFirstCharacter(buf []byte) CrackSingleXorResult {
	return crackSingleXor(buf, nil, histogram.ScoreCaseFirstCharacter)
}

func crackSingleXor(buf []byte, priorHistogram histogram.Histogram, scoreFunc func(histogram.Histogram) float64) CrackSingleXorResult {
	var result CrackSingleXorResult
	for i := range 256 {
		guess := make([]byte, len(buf))
		cipherx.XorByte(guess, buf, byte(i))

		hist := histogram.ComputeHistogram(guess)
		if hist == nil {
			continue
		}
		combinedHist := hist
		if priorHistogram != nil {
			combinedHist.AddVec(hist, priorHistogram)
		}
		score := scoreFunc(combinedHist)
		if score > result.Score {
			result = CrackSingleXorResult{Score: score, Guess: guess, Key: byte(i), GuessHistogram: hist}
		}
	}
	return result
}

type KeySize struct {
	Size         int
	EditDistance float32
}

func GuessXorKeySizes(buf []byte, min int, max int) []KeySize {
	var keySizes []KeySize
	for keySize := min; keySize <= max; keySize++ {
		dist := cipherx.EditDistance(buf[0:keySize], buf[keySize:2*keySize])
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
	for offset := range keySize {
		newBuf := make([]byte, 0, len(buf)/keySize+1)
		for i := offset; i < len(buf); i += keySize {
			newBuf = append(newBuf, buf[i])
		}
		result := CrackSingleXor(newBuf)
		key = append(key, result.Key)
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
		hist := histogram.ComputeHistogram(guess)
		if hist == nil {
			continue
		}
		score := histogram.Score(hist)
		if score > bestScore {
			bestGuess = guess
			bestScore = score
		}
	}

	return bestGuess
}
