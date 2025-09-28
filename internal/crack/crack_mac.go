package crack

import (
	"cryptopals/internal/hashx"
	"cryptopals/internal/hashx/md4x"
	"cryptopals/internal/hashx/sha1x"
	"slices"
)

func ExtendSha1Mac(mac []byte, keySize int, originalMessage []byte, newMessage []byte) ([]byte, []byte) {
	return ExtendMac(sha1x.Utils, mac, keySize, originalMessage, newMessage)
}

func ExtendMd4Mac(mac []byte, keySize int, originalMessage []byte, newMessage []byte) ([]byte, []byte) {
	return ExtendMac(md4x.Utils, mac, keySize, originalMessage, newMessage)
}

// ExtendMac preforms a hash length extension attack
// Constructs the message originalMessage || gluePadding || newMessage
// and returns its respective mac
// Returns (craftedMessage, extendedMac)
func ExtendMac(hu hashx.HashUtils, mac []byte, keySize int, originalMessage []byte, newMessage []byte) ([]byte, []byte) {
	macStateLen := uint64(keySize + len(originalMessage))
	gluePadding := hu.Padding(macStateLen)

	h := hu.FromRegisters(hu.DigestToRegisters(mac), macStateLen+uint64(len(gluePadding)))
	h.Write(newMessage)
	extendedMac := h.Sum(nil) // nil here just means we append the digest to an empty slice

	fullMessage := slices.Concat(originalMessage, gluePadding, newMessage)

	return fullMessage, extendedMac
}
