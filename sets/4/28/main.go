package main

import (
	"cryptopals/internal/cipherx"
	"fmt"
)

func main() {
	msg := []byte("hello world")

	m1 := cipherx.NewSha1Mac([]byte("YELLOW SUBMARINE"))
	sum1 := m1.Sign(msg)

	m2 := cipherx.NewSha1Mac([]byte("YELLOW SUBMARINE 2"))
	sum2 := m2.Sign(msg)

	fmt.Printf("Message signed with key 1: %x\n", sum1)
	fmt.Printf("Message signed with key 2: %x\n", sum2)
}
