package main

import (
	"cryptopals/internal/cipherx"
	"cryptopals/internal/util"
	"fmt"
)

func main() {
	bytes := util.ReadBase64File("10.txt")
	output, _ := cipherx.DecryptAesCbc(bytes, []byte("YELLOW SUBMARINE"), make([]byte, 16))
	fmt.Printf("%s\n", output)
}
