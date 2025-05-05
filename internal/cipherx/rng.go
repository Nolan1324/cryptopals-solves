package cipherx

const n = 624
const m = 397
const w = 32
const r = 31
const uMask = uint32((0xffffffff << r) & 0xffffffff)
const lMask = uint32(0xffffffff) >> (w - r)
const a = uint32(0x9908b0df)
const u = 11
const s = 7
const t = 15
const l = 18
const b = uint32(0x9d2c5680)
const c = uint32(0xefc60000)
const f = uint32(1812433253)

type MersenneTwister struct {
	StateArray [n]uint32
	StateIndex int
}

type Rng interface {
	Rand() uint32
}

func NewMersenneTwister(seed uint32) Rng {
	var mt MersenneTwister

	mt.StateArray[0] = seed
	for i := 1; i < n; i++ {
		seed = f*(seed^(seed>>(w-2))) + uint32(i)
		mt.StateArray[i] = seed
	}

	mt.StateIndex = 0

	return &mt
}

func (mt *MersenneTwister) Rand() uint32 {
	k := mt.StateIndex

	j := k - (n - 1)
	if j < 0 {
		j += n
	}

	x := (mt.StateArray[k] & uMask) | (mt.StateArray[j] & lMask)

	xA := x >> 1
	if x&uint32(0x00000001) != 0 {
		xA ^= a
	}

	j = k - (n - m)
	if j < 0 {
		j += n
	}

	x = mt.StateArray[j] ^ xA
	mt.StateArray[k] = x
	k++

	if k >= n {
		k = 0
	}
	mt.StateIndex = k

	y := x ^ (x >> u)
	y = y ^ ((y << s) & b)
	y = y ^ ((y << t) & c)
	z := y ^ (y >> l)

	return z
}
