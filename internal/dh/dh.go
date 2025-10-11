package dh

import (
	crand "crypto/rand"
	"math/big"
)

const (
	pNistStr = "ffffffffffffffffc90fdaa22168c234c4c6628b80dc1cd129024" +
		"e088a67cc74020bbea63b139b22514a08798e3404ddef9519b3cd" +
		"3a431b302b0a6df25f14374fe1356d6d51c245e485b576625e7ec" +
		"6f44c42e9a637ed6b0bff5cb6f406b7edee386bfb5a899fa5ae9f" +
		"24117c4b1fe649286651ece45b3dc2007cb8a163bf0598da48361" +
		"c55d39a69163fa8fd24cf5f83655d23dca3ad961c62f356208552" +
		"bb9ed529077096966d670c354e4abc9804f1746c08ca237327fff"
)

var pNist *big.Int

// DiffeHellman holds the parameters used to create keys for a Diffe-Hellman key exchange.
type DiffeHellman struct {
	p *big.Int
	g *big.Int
}

type PrivateKey = *big.Int
type PublicKey = *big.Int
type SharedKey = *big.Int

func init() {
	initPNist()
}

func initPNist() {
	p, b := new(big.Int).SetString(pNistStr, 16)
	if !b || p == nil {
		panic("unexpected failure creating big Int")
	}
	pNist = p
}

// MakeDiffeHellman configures a Diffe-Hellman key exchange instance with the provided parameters.
func MakeDiffeHellman(p *big.Int, g *big.Int) DiffeHellman {
	if !p.ProbablyPrime(100) {
		panic("p is not prime")
	}
	return DiffeHellman{p: p, g: g}
}

// MakeDiffeHellman configures a Diffe-Hellman key exchange instance with the standard secure NIST parameters.
func MakeNistDiffeHellman() DiffeHellman {
	return DiffeHellman{p: pNist, g: big.NewInt(2)}
}

// PublicKey computes a client's public key from their private key
func (d DiffeHellman) PublicKey(privateKey PrivateKey) PublicKey {
	return new(big.Int).Exp(d.g, privateKey, d.p)
}

// SharedKey computes the shared key from the client's private key and the other party's public key
func (d DiffeHellman) SharedKey(privateKey PrivateKey, remotePublicKey PublicKey) SharedKey {
	return new(big.Int).Exp(remotePublicKey, privateKey, d.p)
}

// Generate a secure random private key for this Diffe-Hellman key exchange
// Returns an error if there is an issue reading from the OS's random number generator.
func (d DiffeHellman) RandomPrivateKey() (PrivateKey, error) {
	return crand.Int(crand.Reader, d.p)
}

// G is generator parameter for the Diffe-Hellman key exchange
func (d DiffeHellman) G() *big.Int {
	return d.g
}

// P is the modulo parameter for the Diffe-Hellman key exchange
func (d DiffeHellman) P() *big.Int {
	return d.p
}
