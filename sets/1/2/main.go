package main

import (
	"cryptopals/internal/cipherx"
	"cryptopals/internal/enc"
	"fmt"
)

func main() {
	buf1 := enc.HexDecode([]byte("1c0111001f010100061a024b53535009181c"))
	buf2 := enc.HexDecode([]byte("686974207468652062756c6c277320657965"))

	output := cipherx.XorBytes(buf1, buf2)

	fmt.Printf("Expected: 746865206b696420646f6e277420706c6179\n")
	fmt.Printf("Computed: %x\n", output)
}
