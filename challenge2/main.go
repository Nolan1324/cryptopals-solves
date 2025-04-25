package main

import (
	"cryptopals/internal/enc"
	"cryptopals/internal/ops"
	"fmt"
	"log"
)

func main() {
	buf1 := enc.HexDecode([]byte("1c0111001f010100061a024b53535009181c"))
	buf2 := enc.HexDecode([]byte("686974207468652062756c6c277320657965"))

	output, err := ops.XorTwoBuffers(buf1, buf2)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%x\n", output)
}
