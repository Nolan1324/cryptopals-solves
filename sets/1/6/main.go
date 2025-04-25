package main

import (
	"cryptopals/internal/crack"
	"cryptopals/internal/util"
	"fmt"
)

func main() {
	buf := util.ReadBase64File("6.txt")

	guess := crack.CrackRepeatingKeyXor(buf, 2, 100, 7)

	fmt.Printf("%s\n", guess)
}
