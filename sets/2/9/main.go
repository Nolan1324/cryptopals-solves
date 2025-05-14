package main

import (
	"cryptopals/internal/cipherx"
	"fmt"
)

func main() {
	output := cipherx.AddPkcs7Padding([]byte("YELLOW SUBMARINE"), 20)
	fmt.Printf("%q\n", output)
}
