package main

import (
	"cryptopals/internal/ops"
	"fmt"
)

func main() {
	output := ops.Pcks7Padding([]byte("YELLOW SUBMARINE"), 20)
	fmt.Printf("%x\n", output)
}
