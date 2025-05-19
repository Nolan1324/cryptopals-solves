package cipherx

import (
	"bytes"
	"cryptopals/internal/hashx/md4x"
	"cryptopals/internal/hashx/sha1x"
	"hash"
)

type Mac interface {
	Sign(message []byte) []byte
}

type mac struct {
	key  []byte
	hash hash.Hash
}

func NewMac(key []byte, hash hash.Hash) Mac {
	return &mac{key: key, hash: hash}
}

func NewSha1Mac(key []byte) Mac {
	return NewMac(key, sha1x.New())
}

func NewMd4Mac(key []byte) Mac {
	return NewMac(key, md4x.New())
}

func (m mac) Sign(message []byte) []byte {
	m.hash.Reset()
	m.hash.Write(append(bytes.Clone(m.key), message...))
	return m.hash.Sum(nil)
}
