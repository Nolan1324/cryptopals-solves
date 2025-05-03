package main

import (
	"cryptopals/internal/cipherx"
	"cryptopals/internal/randx"
	"fmt"
	"log"
)

func main() {
	app := makeApplication()

	// Perform multiple attempts with random user data each time, because sometimes
	// bitflipping the ciphertext block causes the user data to become invalid
	for range 5 {
		buf := make([]byte, 32)
		for i := range buf {
			buf[i] = byte(randx.RandRange(97, 122+1))
		}

		encrypted, err := app.CreateDataEncrypted(string(buf))
		if err != nil {
			log.Fatal(err)
		}

		blockToFlip := encrypted[32:48]
		flip := make([]byte, 16)
		cipherx.XorBytes(flip, buf[16:32], []byte(";admin=true;a=aa"))
		cipherx.XorBytes(blockToFlip, blockToFlip, flip)

		isAdmin, err := app.IsAdmin(encrypted)
		if err != nil {
			fmt.Printf("Error authenticating: %s\nRetrying\n", err)
			continue
		}
		if isAdmin {
			fmt.Println("Authenticated")
		} else {
			fmt.Println("Not authenticated")
		}
		break
	}
}
