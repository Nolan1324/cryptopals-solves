package main

import (
	"cryptopals/internal/crack"
	"fmt"
)

func main() {
	app := MakeApplication()
	data, mac := app.CreateSignedData()

	for keySize := 1; keySize < 128; keySize++ {
		newMessage := []byte(";admin=true")

		fullMessage, extendedMac := crack.ExtendSha1Mac(mac[:], keySize, data, newMessage)

		isAdmin, err := app.IsAdmin(fullMessage, extendedMac)
		if err != nil {
			continue
		}
		if isAdmin {
			fmt.Printf("Authenticated (key size %v)\n", keySize)
			break
		}
	}
}
