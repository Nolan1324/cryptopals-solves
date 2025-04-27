package cipherx

func XorByte(buf []byte, c byte) []byte {
	output := make([]byte, len(buf))
	for i := range buf {
		output[i] = buf[i] ^ c
	}

	return output
}

func XorBytes(buf1 []byte, buf2 []byte) []byte {
	if len(buf1) != len(buf2) {
		panic("buffers differ in length")
	}

	output := make([]byte, len(buf1))
	for i := range buf1 {
		output[i] = buf1[i] ^ buf2[i]
	}

	return output
}

func RepeatingKeyXor(buf []byte, key []byte) []byte {
	output := make([]byte, len(buf))

	k := len(key)
	for i := range buf {
		output[i] = buf[i] ^ key[i%k]
	}

	return output
}
