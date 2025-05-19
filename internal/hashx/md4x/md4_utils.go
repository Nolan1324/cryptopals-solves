package md4x

import (
	"cryptopals/internal/hashx"
	"encoding/binary"
	"hash"
)

type utils struct{}

var Utils utils

func (utils) FromRegisters(s []uint32, len uint64) hash.Hash {
	if len%chunk != 0 {
		panic("len is not a multiple of chunk")
	}

	var d digest
	d.len = len
	d.nx = 0
	d.s = [4]uint32(s)

	return &d
}

func (utils) DigestToRegisters(digest []byte) []uint32 {
	var s [4]uint32
	s[0] = binary.LittleEndian.Uint32(digest[0:])
	s[1] = binary.LittleEndian.Uint32(digest[4:])
	s[2] = binary.LittleEndian.Uint32(digest[8:])
	s[3] = binary.LittleEndian.Uint32(digest[12:])
	return s[:]
}

func (utils) Padding(len uint64) []byte {
	return hashx.MdPadding(len, binary.LittleEndian)
}
