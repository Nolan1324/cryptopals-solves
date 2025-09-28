package histogram

import (
	"bufio"
	"embed"
	"log"
	"strconv"

	"gonum.org/v1/gonum/mat"
)

type Histogram = *mat.VecDense

//go:embed *.txt
var histogramsFolder embed.FS

var expectedHistogram Histogram
var expectedHistogramFirstChars Histogram

// init is the package-level init function that loads the
// pre-computed histograms into memory.
func init() {
	expectedHistogram = loadHistogramFile("histogram.txt")
	expectedHistogramFirstChars = loadHistogramFile("histogram_first_chars.txt")
}

func NewHistogram() Histogram {
	return mat.NewVecDense(128, nil)
}

func loadHistogramFile(histogramFilename string) Histogram {
	file, err := histogramsFolder.Open(histogramFilename)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()

	var data []float64
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		value, err := strconv.ParseInt(line, 0, 0)
		if err != nil {
			log.Fatal(err)
		}
		data = append(data, float64(value))
	}

	if len(data) != 128 {
		log.Fatalf("histogram in %v is wrong size", histogramFilename)
	}

	histogram := mat.NewVecDense(len(data), data)
	histogram.ScaleVec(1/histogram.Norm(1), histogram)

	return histogram
}

func MakeCaseInsensitive(h Histogram) Histogram {
	histNew := NewHistogram()
	histNew.CopyVec(h)
	for c := 'A'; c <= 'Z'; c++ {
		upperIndex := int(c)
		lowerIndex := int(c) + (int('a') - int('A'))

		histNew.SetVec(upperIndex, histNew.AtVec(upperIndex)+histNew.AtVec(lowerIndex))
		histNew.SetVec(lowerIndex, 0)
	}
	return histNew
}

func incrementFrequency(h Histogram, character byte) {
	index := int(character)
	h.SetVec(index, h.AtVec(index)+1)
}

func ComputeHistogram(text []byte) Histogram {
	vec := NewHistogram()

	for _, c := range text {
		if c >= 128 {
			return nil
		}
		incrementFrequency(vec, c)
	}

	return vec
}

func score(observedHistogram Histogram, expectedHistogram Histogram) float64 {
	if observedHistogram == nil {
		return 0
	}

	numSamples := int(mat.Sum(observedHistogram))

	expectedHistogramScaled := NewHistogram()
	expectedHistogramScaled.ScaleVec(float64(numSamples), expectedHistogram)
	expectedHistogramScaled.ScaleVec(1/expectedHistogramScaled.Norm(2), expectedHistogramScaled)

	observedHistogramScaled := NewHistogram()
	observedHistogramScaled.ScaleVec(1/observedHistogram.Norm(2), observedHistogram)

	return mat.Dot(observedHistogramScaled, expectedHistogramScaled)
}

func Score(observedHistogram Histogram) float64 {
	return score(observedHistogram, expectedHistogram)
}

func ScoreCaseFirstCharacter(observedHistogram Histogram) float64 {
	return score(observedHistogram, expectedHistogramFirstChars)
}
