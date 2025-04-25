package main

import (
	"cryptopals/internal/ops"
	"cryptopals/internal/util"
	"fmt"
)

func main() {
	buf := util.ReadBase64File("7.txt")

	output, _ := ops.DecryptAesEcb(buf, []byte("YELLOW SUBMARINE"))

	fmt.Printf("%s\n", output)
}
