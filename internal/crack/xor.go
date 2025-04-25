package crack

import (
	"cryptopals/internal/ops"
	"cryptopals/internal/util"
	"unicode"

	"gonum.org/v1/gonum/mat"
)

var frequencyVec *mat.VecDense

func initFrequencyVec() {
	if frequencyVec != nil {
		return
	}
	var frequencies = [26]float64{0.082, 0.015, 0.028, 0.043, 0.127, 0.022, 0.02, 0.061, 0.07, 0.0015, 0.0077, 0.04, 0.024, 0.067, 0.075, 0.019, 0.00095, 0.06, 0.063, 0.091, 0.028, 0.0098, 0.024, 0.0015, 0.02, 0.00074}
	f := mat.NewVecDense(26, frequencies[:])
	f.ScaleVec(1/f.Norm(2), f)
	frequencyVec = f
}

func FrequencyScore(guess []byte) float64 {
	initFrequencyVec()

	vec := mat.NewVecDense(26, make([]float64, 26))

	for _, c := range guess {
		if !util.IsWordChar(c) {
			return 0
		}
		if !util.IsAlpha(c) {
			continue
		}
		index := int(unicode.ToLower(rune(c)) - 'a')
		vec.SetVec(index, vec.AtVec(index)+1)
	}

	if vec.Norm(1) == 0 {
		return 0
	}

	vec.ScaleVec(1/vec.Norm(1), vec)
	vec.ScaleVec(1/vec.Norm(2), vec)

	return mat.Dot(frequencyVec, vec)
}

func CrackSingleXor(buf []byte) ([]byte, float64) {
	var best_guess []byte
	var best_score float64
	for i := 0; i < 256; i++ {
		guess := ops.XorByte(buf, byte(i))

		score := FrequencyScore(guess)
		if score > best_score {
			best_score = score
			best_guess = guess
		}
	}
	return best_guess, best_score
}
