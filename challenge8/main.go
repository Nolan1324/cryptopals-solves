package main

import (
	"cryptopals/internal/util"
	"fmt"
)

func score(bytes []byte, bs int) int {
	block_map := make(map[string]int)
	for i := 0; i < len(bytes); i += bs {
		block := string(bytes[i : i+bs])
		block_map[block]++
	}
	score := 0
	for _, v := range block_map {
		if v > 1 {
			score += v
		}
	}
	return score
}

func main() {
	strings := util.ReadHexListFile("8.txt")

	for _, bytes := range strings {
		score := score(bytes, 16)
		if score > 0 {
			fmt.Printf("%x\n", bytes)
		}
	}
}
