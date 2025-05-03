package main

import (
	"testing"
)

func TestAppNormalFunction(t *testing.T) {
	app := makeApplication()
	ct, iv := app.CreateCiphertext()
	if !app.HasValidPadding(ct, iv) {
		t.Fatal("Decrypted ciphertext does not have valid padding")
	}
}
