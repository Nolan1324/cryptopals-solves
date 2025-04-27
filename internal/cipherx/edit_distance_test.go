package cipherx

import (
	"testing"
)

func TestEditDistance(t *testing.T) {
	dist, err := EditDistance([]byte("this is a test"), []byte("wokka wokka!!!"))
	if err != nil {
		t.Fatalf("%v", err)
	}
	if dist != 37 {
		t.Errorf("expected 37, got: %v", dist)
	}
}
