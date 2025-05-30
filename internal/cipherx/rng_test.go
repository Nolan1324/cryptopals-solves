package cipherx

import "testing"

func TestRng(t *testing.T) {
	rng := NewMersenneTwister(5489)
	// https://oeis.org/A221557
	expected := [...]uint32{3499211612, 581869302, 3890346734, 3586334585, 545404204, 4161255391, 3922919429, 949333985, 2715962298, 1323567403, 418932835, 2350294565, 1196140740, 809094426, 2348838239, 4264392720, 4112460519, 4279768804, 4144164697, 4156218106, 676943009, 3117454609}
	for i, expectedVal := range expected {
		val := rng.Rand()
		if val != expectedVal {
			t.Errorf("rng value %v does not match expected value %v at iteration %v", val, expectedVal, i)
		}
	}
}
