package main

import (
	"cryptopals/internal/cipherx"
	"cryptopals/internal/randx"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

type Application struct {
	key []byte
	iv  []byte
}

func makeApplication() Application {
	return Application{key: randx.RandBytes(16), iv: make([]byte, 16)}
}

func (a Application) CreateData(userData string) string {
	return "comment1=cooking%20MCs;userdata=" + url.QueryEscape(userData) + ";comment2=%20like%20a%20pound%20of%20bacon"
}

func (a Application) CreateDataEncrypted(userData string) ([]byte, error) {
	data := cipherx.AddPcks7Padding([]byte(a.CreateData(userData)), 16)
	ct, err := cipherx.EncryptAesCbc(data, a.key, a.iv)
	if err != nil {
		return nil, err
	}
	return ct, nil
}

func (a Application) IsAdmin(encryptedData []byte) (bool, error) {
	data, err := cipherx.DecryptAesCbc(encryptedData, a.key, a.iv)
	if err != nil {
		return false, fmt.Errorf("error decrypting data: %w", err)
	}
	data, err = cipherx.RemovePcks7Padding(data)
	if err != nil {
		return false, fmt.Errorf("error removing decrypted data padding: %w", err)
	}
	// log.Printf("Decrypted data: %q", data)
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
