package histogram

import (
	"math"
	"testing"
)

func TestHistogram(t *testing.T) {
	if expectedHistogram.Len() != 128 {
		t.Fatalf("expected histgram vector length 128, actual %v", expectedHistogram.Len())
	}
	if math.Abs(0.09527-expectedHistogram.AtVec('e')) > 0.001 {
		t.Fatalf("frequency of 'e' is incorrect")
	}
}

func TestScore(t *testing.T) {
	textGood := []byte("hello world")
	textBad := []byte("eeeeeeeeeeeeeeeeeeee")
	scoreGood := Score(ComputeHistogram(textGood))
	scoreBad := Score(ComputeHistogram(textBad))
	t.Logf("Good score: %v", scoreGood)
	t.Logf("Bad score: %v", scoreBad)
	if scoreGood < scoreBad {
		t.Fatalf("good score is less than bad score")
	}
}

func TestScoreCaseFirstCharacter(t *testing.T) {
	text := []byte("HELLO WORLD")
	hist := ComputeHistogram(text)
	score := Score(hist)
	scoreNew := ScoreCaseFirstCharacter(hist)
	t.Logf("Original score: %v", score)
	t.Logf("New score: %v", scoreNew)
	if scoreNew <= score {
		t.Fatalf("new score is <= original score")
	}
}
