package main

import (
	"cryptopals/internal/cipherx"
	"fmt"
	"strconv"
)

func main() {
	output := cipherx.Pcks7Padding([]byte("YELLOW SUBMARINE"), 20)
	fmt.Printf("%s\n", strconv.Quote(string(output)))
}
