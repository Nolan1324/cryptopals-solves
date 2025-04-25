package main

import (
	"bufio"
	"cryptopals/internal/crack"
	"cryptopals/internal/histogram"
	"cryptopals/internal/ops"
	"encoding/base64"
	"fmt"
	"log"
	"os"
)

func loadFile(filename string) []byte {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("error opening file: %v\n", err)
	}
	defer file.Close()

	var buf []byte

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		decoded, err := base64.StdEncoding.DecodeString(line)
		if err != nil {
			log.Fatalf("error \n")
		}
		buf = append(buf, decoded...)
	}

	return buf
}

func main() {
	buf := loadFile("6.txt")
	keySizes := crack.GuessXorKeySizes(buf, 2, 100)
	keySizes = keySizes[:10]

	var bestGuess []byte
	var bestScore float64
	for _, keySize := range keySizes {
		key := crack.CrackRepeatingKeyXor(buf, keySize.Size)
		guess := ops.RepeatingKeyXor(buf, key)
		score := histogram.Score(guess)
		// fmt.Printf("%v\n", histogram.Score(guess))
		if score > bestScore {
			bestGuess = guess
			bestScore = score
		}
	}

	fmt.Printf("%s\n", bestGuess)
}
