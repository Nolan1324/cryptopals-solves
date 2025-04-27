package main

import (
	"cryptopals/internal/cipherx"
	"cryptopals/internal/util"
	"fmt"
)

func main() {
	buf := util.ReadBase64File("7.txt")

	output, _ := cipherx.DecryptAesEcb(buf, []byte("YELLOW SUBMARINE"))

	fmt.Printf("%s\n", output)
}
