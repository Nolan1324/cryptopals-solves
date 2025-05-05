package main

import (
	"cryptopals/internal/cipherx"
	"cryptopals/internal/randx"
	"errors"
	"fmt"
	"time"
)

func generateRandomValue() uint32 {
	sleep := func() {
		time.Sleep(time.Duration(randx.RandRange(10, 60)) * time.Second)
	}
	sleep()
	rng := cipherx.NewMersenneTwister(uint32(time.Now().Unix()))
	sleep()
	return rng.Rand()
}

func crackSeed(randomValue uint32, limit int) (uint32, error) {
	now := uint32(time.Now().Unix())
	for i := range limit {
		seed := now - uint32(i)
		rng := cipherx.NewMersenneTwister(seed)
		if rng.Rand() == randomValue {
			return seed, nil
		}
	}
	return 0, errors.New("no seed found")
}

func main() {
	val := generateRandomValue()
	fmt.Printf("Generated random value: %v\n", val)
	seed, err := crackSeed(val, 3000)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Seed found: %v\n", seed)
}
