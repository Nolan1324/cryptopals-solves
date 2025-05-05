package crack

import (
	"cryptopals/internal/cipherx"
	"math/rand"
	"testing"
)

func TestInverseTemper(t *testing.T) {
	r := rand.New(rand.NewSource(0))
	val := r.Uint32()
	tempered := Temper(val)
	untempered := InverseTemper(tempered)
	if val != untempered {
		t.Errorf("untempered value %x does not match original value %x", untempered, val)
	}
}

func TestInverseLeftShiftAndXor(t *testing.T) {
	r := rand.New(rand.NewSource(0))
	x := r.Uint32()
	k := 7
	m := uint32(0x9d2c5680)
	// m := uint32(0xffffffff)
	y := x ^ ((x << k) & m)
	t.Logf("x = %x, k = %v, m = %x, y = %x", x, k, m, y)
	xNew := inverseLeftShiftAndXor(y, k, m)
	if xNew != x {
		t.Errorf("recovered value %x does not match original value %x", xNew, x)
	}
}

func TestInverseRightShiftAndXor(t *testing.T) {
	r := rand.New(rand.NewSource(0))
	x := r.Uint32()
	k := 11
	y := x ^ (x >> k)
	t.Logf("x = %x, k = %v, y = %x", x, k, y)
	xNew := inverseRightShiftAndXor(y, k)
	if xNew != x {
		t.Errorf("recovered value %x does not match original value %x", xNew, x)
	}
}

func TestCloneRng(t *testing.T) {
	rng := cipherx.NewMersenneTwister(8085)
	output := make([]uint32, 624)
	for i := range output {
		output[i] = rng.Rand()
	}

	rngCloned := CloneRngFromOutput(output)

	for i := range 700 {
		if rng.Rand() != rngCloned.Rand() {
			t.Errorf("Original RNG and cloned RNG differ at value %v\n", i)
		}
	}
}
