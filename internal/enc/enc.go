package enc

import (
	"encoding/base64"
	"encoding/hex"
	"log"
)

func HexToBase64(hex_bytes []byte) []byte {
	return Base64Encode(HexDecode(hex_bytes))
}

func HexDecode(hex_bytes []byte) []byte {
	decoded_len := hex.DecodedLen(len(hex_bytes))
	decoded := make([]byte, decoded_len)
	written, err := hex.Decode(decoded, hex_bytes)
	if err != nil {
		log.Fatal(err)
	}
	if written < decoded_len {
		log.Fatal("Not all bytes written when decoding")
	}

	return decoded
}

func Base64Encode(base64_bytes []byte) []byte {
	encoded_len := base64.StdEncoding.EncodedLen(len(base64_bytes))
	encoded := make([]byte, encoded_len)
	base64.StdEncoding.Encode(encoded, base64_bytes)

	return encoded
}
