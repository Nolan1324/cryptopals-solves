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

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		bytes := enc.HexDecode([]byte(line))
		output, score := crack.CrackSingleXor(bytes)
		if score > 0.6 {
			fmt.Printf("%s %v\n", output, score)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("error reading file: %v", err)
	}
}
