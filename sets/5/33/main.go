package main

import (
	"cryptopals/internal/dh"
	"fmt"
	"math/big"
)

func demoDiffeHellman(diffeHellman dh.DiffeHellman) {
	fmt.Printf("p: %v\ng: %v\n", diffeHellman.P(), diffeHellman.G())
	c1 := dh.MakeClientWithRandomKey(diffeHellman)
	c2 := dh.MakeClientWithRandomKey(diffeHellman)
	fmt.Printf("Client 1 computes public key: %v\n", c1.SharedKey(c2.PublicKey()))
	fmt.Printf("Client 2 computes public key: %v\n", c2.SharedKey(c1.PublicKey()))
}

func main() {
	demoDiffeHellman(dh.MakeDiffeHellman(big.NewInt(37), big.NewInt(2)))
	fmt.Println()
	demoDiffeHellman(dh.MakeNistDiffeHellman())
}
