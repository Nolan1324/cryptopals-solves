package histogram

import (
	_ "embed"
	"log"
	"strconv"
	"strings"

	"gonum.org/v1/gonum/mat"
)

//go:embed histogram.txt
var histogramBytes []byte

var histogramVec *mat.VecDense

func initHistogram() {
	if histogramVec != nil {
		return
	}

	var data []float64
	for _, line := range strings.Split(string(histogramBytes), "\n") {
		value, err := strconv.ParseInt(line, 0, 0)
		if err != nil {
			log.Fatal(err)
		}
		data = append(data, float64(value))
	}

	histogramVec = mat.NewVecDense(len(data), data)
	histogramVec.ScaleVec(1/histogramVec.Norm(1), histogramVec)
}

func Score(guess []byte) float64 {
	initHistogram()

	vec := mat.NewVecDense(128, make([]float64, 128))

	numLetters := 0
	for _, c := range guess {
		if c >= 128 {
			return 0
		}
		index := int(c)
		vec.SetVec(index, vec.AtVec(index)+1)
		numLetters++
	}

	var histogramVecScaled mat.VecDense
	histogramVecScaled.ScaleVec(float64(numLetters), histogramVec)
	histogramVecScaled.ScaleVec(1/histogramVecScaled.Norm(2), &histogramVecScaled)

	vec.ScaleVec(1/vec.Norm(2), vec)

	return mat.Dot(vec, &histogramVecScaled)

	// vec.ScaleVec(1/vec.Norm(1), vec)
	// vec.SubVec(histogramVec, vec)
	// return 1 / vec.Norm(1)
}
