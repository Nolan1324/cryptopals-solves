package main

import (
	"cryptopals/internal/cipherx"
	"fmt"
)

func main() {
	padded := []byte("ICE ICE BABY\x04\x04\x04\x04")
	unpadded, err := cipherx.RemovePcks7Padding(padded)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Remove padding on %q returns %q\n", padded, unpadded)

	padded = []byte("ICE ICE BABY\x05\x05\x05\x05")
	_, err = cipherx.RemovePcks7Padding(padded)
	fmt.Printf("Remove padding on %q returns error '%v'\n", padded, err)

	padded = []byte("ICE ICE BABY\x01\x02\x03\x04")
	_, err = cipherx.RemovePcks7Padding(padded)
	fmt.Printf("Remove padding on %q returns error '%v'\n", padded, err)
}
