package main

import (
	"cryptopals/internal/ops"
	"fmt"
	"strconv"
)

func main() {
	output := ops.Pcks7Padding([]byte("YELLOW SUBMARINE"), 20)
	fmt.Printf("%s\n", strconv.Quote(string(output)))
}
