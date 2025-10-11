package dh

import (
	"math/big"
	"testing"
)

func TestExchange(t *testing.T) {
	diffeHellman := MakeDiffeHellman(big.NewInt(37), big.NewInt(2))
	c1 := MakeClient(diffeHellman, big.NewInt(5))
	c2 := MakeClient(diffeHellman, big.NewInt(14))
	sharedKey1 := c1.SharedKey(c2.PublicKey())
	sharedKey2 := c2.SharedKey(c1.PublicKey())
	if sharedKey1.Cmp(sharedKey2) != 0 {
		t.Fatalf("shared keys do not match")
	}
}
