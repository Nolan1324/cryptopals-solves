package histogram

import (
	"math"
	"testing"
)

func TestHistogram(t *testing.T) {
	initHistogram()
	if histogramVec.Len() != 128 {
		t.Fatalf("expected histgram vector length 128, actual %v", histogramVec.Len())
	}
	if math.Abs(0.09527-histogramVec.AtVec('e')) > 0.001 {
		t.Fatalf("frequency of 'e' is incorrect")
	}
}

func TestScore(t *testing.T) {
	score_good := Score([]byte("hello world"))
	score_bad := Score([]byte("eeeeeeeeeeeeeeeeeeee"))
	t.Logf("Good score: %v", score_good)
	t.Logf("Bad score: %v", score_bad)
	if score_good < score_bad {
		t.Fatalf("good score is less than bad score")
	}
}
