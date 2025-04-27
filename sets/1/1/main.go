package main

import (
	"cryptopals/internal/enc"
	"fmt"
)

func main() {
	input := "49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d"
	output := enc.HexToBase64([]byte(input))
	fmt.Printf("Expected: SSdtIGtpbGxpbmcgeW91ciBicmFpbiBsaWtlIGEgcG9pc29ub3VzIG11c2hyb29t\n")
	fmt.Printf("Computed: %s\n", output)
}
