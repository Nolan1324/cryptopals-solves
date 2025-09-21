package main

import (
	"cryptopals/internal/apps/timeattack"
	"encoding/hex"
	"fmt"
	"time"
)

const endpoint = "http://localhost:9000/test"

func attack(file string, numSamples int) ([]byte, error) {
	guess := make([]byte, 20)

	for i := range guess {
		var bestDuration time.Duration
		var bestByte byte
		for b := range 256 {
			var meanDuration time.Duration
			for range numSamples {
				guess[i] = byte(b)
				duration, ok, err := timeattack.DoRequest(endpoint, file, hex.EncodeToString(guess))
				if err != nil {
					return nil, err
				}
				if ok {
					return guess, nil
				}
				meanDuration += duration
			}
			meanDuration /= time.Duration(meanDuration.Seconds() / float64(numSamples))

			if meanDuration > bestDuration {
				bestDuration, bestByte = meanDuration, byte(b)
			}
		}
		guess[i] = bestByte
		fmt.Printf("Current guess: %x\n", guess[:i+1])
	}

	return nil, nil
}

func main() {
	// cancel := timeattack.RunServerProcessInBackground("./server")
	// defer cancel()
	// time.Sleep(3 * time.Second)

	guess, err := attack("foo", 10)
	if err != nil {
		fmt.Printf("error attacking: %v\n", err)
	}
	if guess != nil {
		fmt.Printf("Signature found: %x\n", guess)
	} else {
		fmt.Printf("No signature found\n")
	}
}
