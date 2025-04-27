package ops

import (
	"crypto/aes"
)

func EncryptAesEcb(buf []byte, key []byte) ([]byte, error) {
	cipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	encrypted := make([]byte, len(buf))
	bs := cipher.BlockSize()
	for i := 0; i < len(buf); i += bs {
		cipher.Encrypt(encrypted[i:i+bs], buf[i:i+bs])
	}

	return encrypted, nil
}

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

func EncryptAesCbc(buf []byte, key []byte, iv []byte) ([]byte, error) {
	cipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	prevCiphertextBlock := iv
	encrypted := make([]byte, len(buf))
	bs := cipher.BlockSize()
	for i := 0; i < len(buf); i += bs {
		block := XorBytes(prevCiphertextBlock, buf[i:i+bs])
		cipher.Encrypt(encrypted[i:i+bs], block)
		prevCiphertextBlock = encrypted[i : i+bs]
	}

	return encrypted, nil
}

func DecryptAesCbc(buf []byte, key []byte, iv []byte) ([]byte, error) {
	cipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	prevCiphertextBlock := iv
	decrypted := make([]byte, len(buf))
	bs := cipher.BlockSize()
	for i := 0; i < len(buf); i += bs {
		cipher.Decrypt(decrypted[i:i+bs], buf[i:i+bs])
		block := XorBytes(prevCiphertextBlock, decrypted[i:i+bs])
		copy(decrypted[i:i+bs], block)
		prevCiphertextBlock = buf[i : i+bs]
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
