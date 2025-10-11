package dh

import "math/big"

// DiffeHellmanClient is a client of a Diffe-Hellman key exchange that holds its own private and public key
type DiffeHellmanClient struct {
	diffeHellman DiffeHellman
	privateKey   PrivateKey
	publicKey    PublicKey
}

// DiffeHellmanClient makes a clent for a Diffe-Hellman key exchange, given the parameters and their private key
func MakeClient(diffeHellman DiffeHellman, privateKey PrivateKey) DiffeHellmanClient {
	return DiffeHellmanClient{
		diffeHellman: diffeHellman,
		privateKey:   privateKey,
		publicKey:    diffeHellman.PublicKey(privateKey),
	}
}

// DiffeHellmanClient makes a clent for a Diffe-Hellman key exchange with a random private key, given the parameters
func MakeClientWithRandomKey(diffeHellman DiffeHellman) DiffeHellmanClient {
	privateKey := diffeHellman.RandomPrivateKey()
	return MakeClient(diffeHellman, privateKey)
}

// PublicKey is the client's public key to send to the other party
func (c DiffeHellmanClient) PublicKey() PublicKey {
	return c.publicKey
}

// SharedKey computes the shared key from the other parties public key
func (c DiffeHellmanClient) SharedKey(remotePublicKey PublicKey) SharedKey {
	return c.diffeHellman.SharedKey(c.privateKey, remotePublicKey)
}

// G is generator parameter
func (c DiffeHellmanClient) G() *big.Int {
	return c.diffeHellman.G()
}

// P is generator parameter
func (c DiffeHellmanClient) P() *big.Int {
	return c.diffeHellman.P()
}
