package main

import (
	"cryptopals/internal/crack"
	"cryptopals/internal/enc"
	"fmt"
)

func main() {
	buf := enc.HexDecode([]byte("1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736"))

	result := crack.CrackSingleXor(buf)
	fmt.Printf("%s %v %v\n", result.Guess, result.Key, result.Score)
}
