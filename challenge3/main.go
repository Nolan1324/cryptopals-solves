package main

import (
	"cryptopals/internal/crack"
	"cryptopals/internal/enc"
	"cryptopals/internal/util"
	"fmt"
)

func main() {
	buf := enc.HexDecode([]byte("1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736"))

	guess, score := crack.CrackSingleXor(buf, util.IsSentenceAscii)
	fmt.Printf("%s %v\n", guess, score)
}
