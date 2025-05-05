package main

import (
	"cryptopals/internal/cipherx"
	"cryptopals/internal/crack"
	"fmt"
	"os"
)

func main() {
	rng := cipherx.NewMersenneTwister(8085)
	output := make([]uint32, 624)
	for i := range output {
		output[i] = rng.Rand()
	}

	rngCloned := crack.CloneRngFromOutput(output)

	fmt.Printf("Original RNG returns %v on next value\n", rng.Rand())
	fmt.Printf("Cloned RNG returns %v on next value\n", rngCloned.Rand())

	for range 700 {
		if rng.Rand() != rngCloned.Rand() {
			fmt.Printf("Original RNG and cloned RNG differ at some value\n")
			os.Exit(1)
		}
	}
	fmt.Printf("Original RNG and cloned RNG are identical for the next 700 values\n")
}
