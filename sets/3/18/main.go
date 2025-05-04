package main

import (
	"cryptopals/internal/cipherx"
	"encoding/base64"
	"fmt"
)

func main() {
	ct, err := base64.StdEncoding.DecodeString("L77na/nrFsKvynd6HzOoG7GHTLXsTVu9qvY/2syLXzhPweyyMTJULu/6/kXX0KSvoOLSFQ==")
	if err != nil {
		panic(err)
	}
	cipher, err := cipherx.NewAesCtr([]byte("YELLOW SUBMARINE"))
	if err != nil {
		panic(err)
	}
	pt := cipher.Decrypt(ct, 0)
	fmt.Printf("%q\n", pt)
}
