package dh

import (
	"math/big"
	"testing"
)

func TestPublicKey(t *testing.T) {
	diffeHellman := MakeDiffeHellman(big.NewInt(37), big.NewInt(2))
	privKey := big.NewInt(6)
	// 2^6 = 64. 64 % 37 = 27
	expectedPubKey := big.NewInt(27)

	pubKey := diffeHellman.PublicKey(privKey)

	if pubKey.Cmp(expectedPubKey) != 0 {
		t.Fatalf("expected public key %v, got %v", expectedPubKey, pubKey)
	}
}

func TestShared(t *testing.T) {
	diffeHellman := MakeDiffeHellman(big.NewInt(37), big.NewInt(2))
	privKey := big.NewInt(4)
	remotePubKey := big.NewInt(3)
	// 3^4 = 81. 81 % 37 = 7
	expectedSharedKey := big.NewInt(7)

	sharedKey := diffeHellman.SharedKey(privKey, remotePubKey)

	if sharedKey.Cmp(expectedSharedKey) != 0 {
		t.Fatalf("expected shared key %v, got %v", expectedSharedKey, sharedKey)
	}
}

func TestNistConstructor(t *testing.T) {
	// ensure no panic occurs
	MakeNistDiffeHellman()
}
