package main

import (
	"cryptopals/internal/ops"
	"fmt"
)

func main() {
	buf := []byte("Burning 'em, if you ain't quick and nimble\nI go crazy when I hear a cymbal")
	output := ops.RepeatingKeyXor(buf, []byte("ICE"))
	fmt.Printf("%x\n", output)
}
