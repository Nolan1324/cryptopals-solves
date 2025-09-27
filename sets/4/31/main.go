package main

import (
	"cryptopals/internal/apps/timeattack"
	"cryptopals/internal/randx"
	"encoding/hex"
	"log"
	"net/http"
	"time"
)

const (
	address  = "localhost:8000"
	endpoint = "http://" + address + "/" + timeattack.TestEndpoint
	hmacLen  = 20
)

func attack(file string) ([]byte, error) {
	guess := make([]byte, hmacLen)

	for i := range guess {
		var bestDuration time.Duration
		var bestByte byte
		for b := range 256 {
			guess[i] = byte(b)
			duration, ok, err := timeattack.DoRequest(endpoint, file, hex.EncodeToString(guess))
			if err != nil {
				return nil, err
			}
			if ok {
				return guess, nil
			}
			if duration > bestDuration {
				bestDuration, bestByte = duration, byte(b)
			}
		}
		guess[i] = bestByte
		log.Printf("Current guess: %x\n", guess[:i+1])
	}

	return nil, nil
}

func main() {
	server := timeattack.NewServer(address, randx.RandBytes(16), time.Millisecond*5, false)
	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
	defer func() {
		err := server.Close()
		if err != nil {
			log.Printf("when closing server: %v\n", err)
		}
	}()

	file := "foo"
	answer := server.Sign([]byte(file))
	log.Printf("True answer: %x", answer)

	err := timeattack.WaitForServerStartWithTimeout(3*time.Second, 5*time.Millisecond, endpoint)
	if err != nil {
		log.Fatalf("timeout out waiting for server: %v", err)
	}

	guess, err := attack(file)
	if err != nil {
		log.Fatalf("when attacking: %v\n", err)
	}
	if guess != nil {
		log.Printf("Signature for '%v' found: %x\n", file, guess)
	} else {
		log.Printf("No signature for '%v' found\n", file)
	}
}
