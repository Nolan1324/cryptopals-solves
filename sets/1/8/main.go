package main

import (
	"cryptopals/internal/crack"
	"cryptopals/internal/util"
	"fmt"
)

func main() {
	strings := util.ReadHexListFile("8.txt")

	for _, bytes := range strings {
		if crack.HasRepeatingBlock(bytes, 16) {
			fmt.Printf("%x\n", bytes)
		}
	}
}
