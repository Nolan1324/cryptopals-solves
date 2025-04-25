package main

import (
	"cryptopals/internal/crack"
	"cryptopals/internal/util"
	"fmt"
)

func main() {
	strings := util.ReadHexListFile("4.txt")

	var bestGuess []byte
	var bestScore float64

	for _, bytes := range strings {
		guess, _, score := crack.CrackSingleXor(bytes)
		if score > bestScore {
			bestScore = score
			bestGuess = guess
		}
	}

	fmt.Printf("%s\n", bestGuess)
}
