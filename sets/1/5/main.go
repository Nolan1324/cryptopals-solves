package main

import (
	"cryptopals/internal/cipherx"
	"fmt"
)

func main() {
	buf := []byte("Burning 'em, if you ain't quick and nimble\nI go crazy when I hear a cymbal")
	output := make([]byte, len(buf))
	cipherx.RepeatingKeyXor(output, buf, []byte("ICE"))
	fmt.Printf("Expected: 0b3637272a2b2e63622c2e69692a23693a2a3c6324202d623d63343c2a26226324272765272a282b2f20430a652e2c652a3124333a653e2b2027630c692b20283165286326302e27282f\n")
	fmt.Printf("Computed: %x\n", output)
}
