package hashx

import (
	"hash"
)

type HashUtils interface {
	// FromRegisters reconstructs the hash state from the registers, assuming that the previously summed data was already padded.
	// h is the list of registers, and have exactly the same number of registers as the respective hash algorithm.
	// len is the length of the previously summed data and must be a multiple of 64.
	FromRegisters(h []uint32, len uint64) hash.Hash

	// DigestToRegisters splits the digest into 32-bit registers.
	// digest must be the correct length for the respective hash algorithm.
	DigestToRegisters(digest []byte) []uint32

	// Generates MD-padding for a message of length len for the respective hash algorithm.
	Padding(len uint64) []byte
}
