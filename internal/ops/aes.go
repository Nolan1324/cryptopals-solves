package ops

import (
	"crypto/aes"
)

func DecryptAesEcb(buf []byte, key []byte) ([]byte, error) {
	cipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	decrypted := make([]byte, len(buf))
	bs := cipher.BlockSize()
	for i := 0; i < len(buf); i += bs {
		cipher.Decrypt(decrypted[i:i+bs], buf[i:i+bs])
	}

	return decrypted, nil
}

func Pcks7Padding(buf []byte, bs int) []byte {
	paddingSize := bs - (len(buf) % bs)
	output := make([]byte, len(buf), len(buf)+paddingSize)
	copy(output, buf)
	for i := 0; i < paddingSize; i++ {
		output = append(output, byte(paddingSize))
	}
	return output
}
