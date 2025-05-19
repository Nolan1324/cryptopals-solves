package sha1x

import (
	"cryptopals/internal/hashx"
	"encoding/binary"
	"hash"
)

type utils struct{}

var Utils utils

func (utils) FromRegisters(h []uint32, hashedDataLen uint64) hash.Hash {
	if hashedDataLen%chunk != 0 {
		panic("hashedDataLen is not a multiple of chunk")
	}
	if len(h) != 5 {
		panic("h must be length 5")
	}

	var d digest
	d.len = hashedDataLen
	d.nx = 0
	copy(d.h[:], h)

	return &d
}

func (utils) DigestToRegisters(digest []byte) []uint32 {
	var h [5]uint32
	h[0] = binary.BigEndian.Uint32(digest[0:])
	h[1] = binary.BigEndian.Uint32(digest[4:])
	h[2] = binary.BigEndian.Uint32(digest[8:])
	h[3] = binary.BigEndian.Uint32(digest[12:])
	h[4] = binary.BigEndian.Uint32(digest[16:])
	return h[:]
}

func (utils) Padding(len uint64) []byte {
	return hashx.MdPadding(len, binary.BigEndian)
}
