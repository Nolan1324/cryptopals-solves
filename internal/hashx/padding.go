// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hashx

import "encoding/binary"

func MdPadding(len uint64, byteOrder binary.ByteOrder) []byte {
	// Padding.  Add a 1 bit and 0 bits until 56 bytes mod 64.
	var tmp [64 + 8]byte // padding + length buffer
	tmp[0] = 0x80
	var t uint64
	if len%64 < 56 {
		t = 56 - len%64
	} else {
		t = 64 + 56 - len%64
	}

	// Length in bits.
	len <<= 3
	padlen := tmp[:t+8]
	byteOrder.PutUint64(padlen[t:], len)

	return padlen
}
