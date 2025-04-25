package main

import (
	"bufio"
	"cryptopals/internal/crack"
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

	guess := crack.CrackRepeatingKeyXor(buf, 2, 100, 7)

	fmt.Printf("%s\n", guess)
}
