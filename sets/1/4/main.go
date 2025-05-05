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
		result := crack.CrackSingleXor(bytes)
		if result.Score > bestScore {
			bestScore = result.Score
			bestGuess = result.Guess
		}
	}

	fmt.Printf("%s\n", bestGuess)
}
