package main

import (
	"cryptopals/internal/cipherx"
	"fmt"
)

func main() {
	// Should generate the same values as those listed in https://oeis.org/A221557
	rng := cipherx.NewMersenneTwister(5489)
	for range 22 {
		fmt.Printf("%v ", rng.Rand())
	}
	fmt.Println()
}
