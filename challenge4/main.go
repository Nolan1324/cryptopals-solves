package main

import (
	"bufio"
	"cryptopals/internal/crack"
	"cryptopals/internal/enc"
	"fmt"
	"log"
	"os"
)

func main() {
	file, err := os.Open("4.txt")
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()

	var best_guess []byte
	var best_score float64

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		bytes := enc.HexDecode([]byte(line))
		guess, _, score := crack.CrackSingleXor(bytes)
		if score > best_score {
			best_score = score
			best_guess = guess
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("error reading file: %v", err)
	}

	fmt.Printf("%s\n", best_guess)
}
