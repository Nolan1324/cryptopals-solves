package cipherx

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
		panic("buffers differ in length")
	}

	dist := 0
	for i := range buf1 {
		dist += editDistanceByte(buf1[i], buf2[i])
	}

	return dist, nil
}
