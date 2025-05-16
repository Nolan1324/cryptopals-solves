package main

import (
	"cryptopals/internal/cipherx"
	"cryptopals/internal/slicex"
	"fmt"
	"os"
)

func main() {
	app := MakeApplication()
	const ptLen = 12
	pt := slicex.Repeat(byte('a'), ptLen)
	ct, nonce := app.CreateDataEncrypted(string(pt))

	newPt := []byte("a;admin=true")
	if len(newPt) != ptLen {
		panic("new plaintext length differs from original plaintext")
	}

	keyStream := make([]byte, ptLen)
	cipherx.XorBytes(keyStream, ct[32:32+ptLen], pt)
	cipherx.XorBytes(ct[32:32+ptLen], keyStream, newPt)

	isAdmin, err := app.IsAdmin(ct, nonce)
	if err != nil {
		panic(err)
	}
	if isAdmin {
		fmt.Println("Authenticated")
	} else {
		fmt.Println("Not authenticated")
		os.Exit(1)
	}
}
