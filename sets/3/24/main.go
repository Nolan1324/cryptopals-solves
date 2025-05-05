package main

import (
	"bytes"
	"crypto/rand"
	"cryptopals/internal/cipherx"
	"cryptopals/internal/randx"
	"encoding/binary"
	"fmt"
	"time"
)

func RandUint16() uint16 {
	buf := make([]byte, 2)
	rand.Read(buf)
	return binary.LittleEndian.Uint16(buf)
}

func crackMersenneTwisterCipherSeed(ct []byte, pt []byte) uint16 {
	for seed := range uint32(1 << 16) {
		rng := cipherx.NewMersenneTwister(seed)
		for range len(ct) - len(pt) {
			_ = rng.Rand()
		}
		found := true
		for i, b := range ct[len(ct)-len(pt):] {
			if byte(rng.Rand())^b != pt[i] {
				found = false
				break
			}
		}
		if found {
			return uint16(seed)
		}
	}
	panic("no seed found")
}

func checkPasswordResetToken(token []byte, maxIterations int) bool {
	for iteration := range maxIterations {
		rng := cipherx.NewMersenneTwister(uint32(time.Now().Unix() - int64(iteration)))
		guessToken := make([]byte, 16)
		for i := range guessToken {
			guessToken[i] = byte(rng.Rand())
		}
		if bytes.Equal(token, guessToken) {
			return true
		}
	}
	return false
}

func part1() {
	pt := []byte("AAAA")
	key, ct := encryptWithMersenneTwisterCipher(pt)
	crackedKey := crackMersenneTwisterCipherSeed(ct, pt)
	fmt.Printf("Cracked key %v\n", crackedKey)
	fmt.Printf("True key was %v\n", key)
}

func part2() {
	realToken := generatePasswordResetToken()
	fakeToken := randx.RandBytes(len(realToken))
	time.Sleep(time.Duration(randx.RandRange(4, 10)) * time.Second)
	fmt.Printf("Real token detected as valid=%v\n", checkPasswordResetToken(realToken, 1000))
	fmt.Printf("Fake token detected as valid=%v\n", checkPasswordResetToken(fakeToken, 1000))
}

func main() {
	part1()
	fmt.Println()
	part2()
}
