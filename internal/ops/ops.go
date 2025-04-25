package ops

import "fmt"

func XorSingleByte(buf []byte, c byte) []byte {
	output := make([]byte, len(buf))
	for i := range buf {
		output[i] = buf[i] ^ c
	}

	return output
}

func XorTwoBuffers(buf1 []byte, buf2 []byte) ([]byte, error) {
	if len(buf1) != len(buf2) {
		return nil, fmt.Errorf("buffers differ in length")
	}

	output := make([]byte, len(buf1))
	for i := range buf1 {
		output[i] = buf1[i] ^ buf2[i]
	}

	return output, nil
}

func RepeatingKeyXor(buf []byte, key []byte) []byte {
	output := make([]byte, len(buf))

	k := len(key)
	for i := range buf {
		output[i] = buf[i] ^ key[i%k]
	}

	return output
}

func editDistanceByte(a byte, b byte) int {
	d := a ^ b
	dist := 0
	for d != 0 {
		if d&1 == 1 {
			dist++
		}
		d = d >> 1
	}
	return dist
}

func EditDistance(buf1 []byte, buf2 []byte) (int, error) {
	if len(buf1) != len(buf2) {
		return 0, fmt.Errorf("buffers differ in length")
	}

	dist := 0
	for i := range buf1 {
		dist += editDistanceByte(buf1[i], buf2[i])
	}

	return dist, nil
}
