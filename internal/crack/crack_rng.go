package crack

import "cryptopals/internal/cipherx"

const n = 624

const u = 11
const s = 7
const t = 15
const l = 18
const b = uint32(0x9d2c5680)
const c = uint32(0xefc60000)

func CloneRngFromOutput(output []uint32) cipherx.Rng {
	if len(output) != n {
		panic("length of output must be 624")
	}
	var mt cipherx.MersenneTwister
	for i, val := range output {
		mt.StateArray[i] = InverseTemper(val)
	}
	return &mt
}

func Temper(x uint32) uint32 {
	y := x ^ (x >> u)
	y = y ^ ((y << s) & b)
	y = y ^ ((y << t) & c)
	z := y ^ (y >> l)
	return z
}

func InverseTemper(z uint32) uint32 {
	y := z ^ (z >> l)
	y = y ^ ((y << t) & c)
	y = inverseLeftShiftAndXor(y, s, b)
	x := inverseRightShiftAndXor(y, u)
	return x
}

// inverseLeftShiftAndXor computes the inverse of the operation y = x ^ ((y << k) & m)
func inverseLeftShiftAndXor(y uint32, k int, m uint32) uint32 {
	if k <= 0 || k > 32 {
		panic("invalid k value")
	}
	kMask := 1<<uint32(k) - uint32(1)
	x := y & kMask
	for i := 1; i*k < 32; i++ {
		kMask = kMask << k
		x |= (y & kMask) ^ ((x << k) & kMask & m)
	}
	return x
}

func inverseRightShiftAndXor(y uint32, k int) uint32 {
	if k <= 0 || k > 32 {
		panic("invalid k value")
	}
	kMask := (1<<uint32(k) - uint32(1)) << (32 - k)
	x := y & kMask
	for i := 1; i*k < 32; i++ {
		kMask = kMask >> k
		x |= (y & kMask) ^ ((x >> k) & kMask)
	}
	return x
}
