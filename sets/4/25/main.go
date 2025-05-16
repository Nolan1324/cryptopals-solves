package main

import (
	"cryptopals/internal/cipherx"
	"fmt"
)

func main() {
	app := MakeApplication()
	ct := app.ReadCiphertext()

	for i := range ct {
		app.EditByte(i, 0x00)
	}
	key := app.ReadCiphertext()
	pt := make([]byte, len(ct))
	cipherx.XorBytes(pt, ct, key)
	fmt.Printf("%q\n", pt)
}
