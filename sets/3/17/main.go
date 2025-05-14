package main

import (
	"cryptopals/internal/cipherx"
	"fmt"
	"log"
	"slices"
)

const bs = 16

func crackCbcBlock(prevBlock []byte, block []byte, app Application) []byte {
	iv := slices.Clone(prevBlock)
	decrypted := make([]byte, bs)
	for i := bs - 1; i >= 0; i-- {
		numPadding := bs - i
		var found bool
		// Try tampering iv[i] until we determine which byte
		// causes a valid padding byte at block[i] when decrypted
		for c := range 256 {
			c := byte(c)
			iv[i] = c
			if app.HasValidPadding(block, iv) {
				// Sometimes tampering the last byte can create a valid padding other than \x01.
				// For instance, this occurs if the last byte was change to \x02 and the penultimate byte happens to be \x02 already.
				// To check for this case, tamper the penultimate byte. If the padding becomes invalid, then we do not truly
				// have \x01 in the last byte, so keep searching.
				if i == bs-1 {
					ivTemp := slices.Clone(iv)
					ivTemp[i-1] ^= 1
					if !app.HasValidPadding(block, ivTemp) {
						continue
					}
				}

				if found {
					log.Fatalf("multiple possibilities were found for plaintext block byte i=%v", i)
				}
				decrypted[i] = prevBlock[i] ^ c ^ byte(numPadding)
				found = true
			}
		}
		if !found {
			log.Fatalf("no solution found for plaintext block byte i=%v", i)
		}
		// Tamper bytes [i, bs) to have the correct padding byte for the next iteration
		for j := i; j < bs; j++ {
			iv[j] = prevBlock[j] ^ decrypted[j] ^ byte(numPadding+1)
		}
	}
	return decrypted
}

func main() {
	app := makeApplication()
	ct, iv := app.CreateCiphertext()
	prevBlock := iv
	var decrypted []byte
	for i := 0; i < len(ct); i += bs {
		block := ct[i : i+bs]
		decryptedBlock := crackCbcBlock(prevBlock, block, app)
		decrypted = append(decrypted, decryptedBlock...)
		prevBlock = block
	}
	decrypted, err := cipherx.RemovePkcs7Padding(decrypted)
	if err != nil {
		fmt.Printf("decrypted plaintext does not have valid padding\n")
	}
	fmt.Printf("%q\n", decrypted)
}
