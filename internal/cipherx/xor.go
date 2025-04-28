package cipherx

func XorByte(dst []byte, src []byte, c byte) {
	if len(dst) < len(src) {
		panic("dst buffer is smaller than src")
	}

	for i := range src {
		dst[i] = src[i] ^ c
	}
}

func XorBytes(dst []byte, src1 []byte, src2 []byte) {
	if len(src1) != len(src2) {
		panic("src buffers differ in length")
	}
	if len(dst) < len(src1) {
		panic("dst buffer is smaller than src")
	}

	for i := range src1 {
		dst[i] = src1[i] ^ src2[i]
	}
}

func RepeatingKeyXor(dst []byte, src []byte, key []byte) {
	if len(dst) < len(src) {
		panic("dst buffer is smaller than src")
	}

	k := len(key)
	for i := range src {
		dst[i] = src[i] ^ key[i%k]
	}
}
