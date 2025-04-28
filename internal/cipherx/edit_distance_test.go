package cipherx

import (
	"testing"
)

func TestEditDistance(t *testing.T) {
	dist := EditDistance([]byte("this is a test"), []byte("wokka wokka!!!"))
	if dist != 37 {
		t.Errorf("expected 37, got: %v", dist)
	}
}
