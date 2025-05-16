package main

import (
	"cryptopals/internal/cipherx"
	"cryptopals/internal/randx"
	"errors"
	"fmt"
	"math/rand/v2"
	"net/url"
	"strings"
)

type Application struct {
	ctr cipherx.AesCtr
}

func MakeApplication() Application {
	ctr, err := cipherx.MakeAesCtr(randx.RandBytes(16))
	if err != nil {
		panic(err)
	}
	return Application{ctr: ctr}
}

func (a Application) CreateData(userData string) string {
	return "comment1=cooking%20MCs;userdata=" + url.QueryEscape(userData) + ";comment2=%20like%20a%20pound%20of%20bacon"
}

// CreateDataEncrypted creates a string containing "userdata=<userData>;" signed with AES CTR.
// It returns the AES CTR ciphertext along with the nonce
func (a Application) CreateDataEncrypted(userData string) ([]byte, cipherx.Count) {
	nonce := rand.Uint64()
	ct := a.ctr.Encrypt([]byte(a.CreateData(userData)), nonce)
	return ct, nonce
}

// IsAdmin takes in encryptedData signed with AES CTR along with the nonce, and checks if the decrypted
// data contains "admin=true;"
func (a Application) IsAdmin(encryptedData []byte, nonce cipherx.Count) (bool, error) {
	data := a.ctr.Decrypt(encryptedData, nonce)
	isAdmin, err := isAdminDecrypted(string(data))
	if err != nil {
		return false, fmt.Errorf("error verifying data: %w", err)
	}
	return isAdmin, nil
}

func isAdminDecrypted(data string) (bool, error) {
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
