// This is almost exactly the same as 29/app.go,
// except it uses NewMd4Mac instead of NewSha1Mac

package main

import (
	"bytes"
	"cryptopals/internal/cipherx"
	"cryptopals/internal/randx"
	"errors"
	"strings"
)

type Application struct {
	mac cipherx.Mac
}

func MakeApplication() Application {
	key := randx.RandBytes(randx.RandRange(8, 64))
	return Application{mac: cipherx.NewMd4Mac(key)}
}

// CreateSignedData creates a fixed data string along with its signature, signed with a private key
// Returns (data, signature)
func (a Application) CreateSignedData() ([]byte, []byte) {
	const data = "comment1=cooking%20MCs;userdata=foo;comment2=%20like%20a%20pound%20of%20bacon"
	dataBytes := []byte(data)
	return dataBytes, a.mac.Sign([]byte(dataBytes))
}

// IsAdmin first verifies that signature was produced by signing data with the private key
// If the signature is invalid, an error is returned.
// Otherwise, it checks if the data has admin=true set.
func (a Application) IsAdmin(data []byte, signature []byte) (bool, error) {
	if !bytes.Equal(a.mac.Sign(data), signature[:]) {
		return false, errors.New("signature does not match data")
	}
	isAdmin, err := isAdmin(string(data))
	if err != nil {
		return false, err
	}
	return isAdmin, nil
}

func isAdmin(data string) (bool, error) {
	for _, s := range strings.Split(data, ";") {
		parts := strings.Split(s, "=")
		if len(parts) != 2 {
			return false, errors.New("invalid data string")
		}
		if parts[0] == "admin" && parts[1] == "true" {
			return true, nil
		}
	}
	return false, nil
}
