package cipherx

import (
	"bytes"
	"testing"
)

func TestEcbIdentity(t *testing.T) {
	plaintext := []byte("this message has forty-eight characters exactly.")
	encrypted, _ := EncryptAesEcb(plaintext, []byte("YELLOW SUBMARINE"))
	decrypted, _ := DecryptAesEcb(encrypted, []byte("YELLOW SUBMARINE"))
	if !bytes.Equal(plaintext, decrypted) {
		t.Fail()
	}
}

func TestPcks7Padding(t *testing.T) {
	output := Pcks7Padding([]byte("YELLOW SUBMARINE"), 20)
	if !bytes.Equal(output, []byte("YELLOW SUBMARINE\x04\x04\x04\x04")) {
		t.Fail()
	}
}

func TestPcks7PaddingAligned(t *testing.T) {
	// When aligned to block, padding should still pad a block of \x10's
	output := Pcks7Padding([]byte("0123456789ABCDEF"), 16)
	if !bytes.Equal(output, []byte("0123456789ABCDEF\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10")) {
		t.Fail()
	}
}

func TestPcks7AddRemove(t *testing.T) {
	text := []byte("YELLOW SUBMARINE")
	padded := Pcks7Padding(text, 20)
	unpadded := RemovePcks7Padding(padded)
	if !bytes.Equal(text, unpadded) {
		t.Errorf("original text %q does not match unpadded text %q", text, unpadded)
	}
}
