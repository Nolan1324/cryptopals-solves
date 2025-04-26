package main

import (
	"cryptopals/internal/ops"
	"cryptopals/internal/util"
	"fmt"
)

func main() {
	bytes := util.ReadBase64File("10.txt")
	output, _ := ops.DecryptAesCbc(bytes, []byte("YELLOW SUBMARINE"), make([]byte, 16))
	fmt.Printf("%s\n", output)
}
