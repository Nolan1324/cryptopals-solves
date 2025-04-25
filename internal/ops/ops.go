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
