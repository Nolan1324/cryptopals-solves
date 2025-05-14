package cipherx

import (
	"crypto/aes"
	"errors"
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
		block := make([]byte, bs)
		XorBytes(block, prevCiphertextBlock, buf[i:i+bs])
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
		block := make([]byte, bs)
		cipher.Decrypt(block, buf[i:i+bs])
		XorBytes(decrypted[i:i+bs], prevCiphertextBlock, block)
		prevCiphertextBlock = buf[i : i+bs]
	}

	return decrypted, nil
}

// AddPcks7Padding function adds PCKS7 padding to the end of a slice. If
// it has sufficient capacity, the destination is resliced to accommodate the
// new elements. If it does not, a new underlying array will be allocated.
// AddPcks7Padding returns the updated slice. It is therefore necessary to store the
// result of AddPcks7Padding, often in the variable holding the slice itself:
//
//	slice = append(slice, 16)
func AddPcks7Padding(buf []byte, bs int) []byte {
	paddingSize := bs - (len(buf) % bs)
	for range paddingSize {
		buf = append(buf, byte(paddingSize))
	}
	return buf
}

func RemovePcks7Padding(buf []byte) ([]byte, error) {
	err := errors.New("invalid padding")
	if len(buf) == 0 {
		return buf, err
	}
	paddingByte := buf[len(buf)-1]
	paddingLen := int(paddingByte)
	if paddingLen < 1 || paddingLen > len(buf) {
		return buf, err
	}
	for _, b := range buf[len(buf)-paddingLen:] {
		if b != paddingByte {
			return buf, err
		}
	}
	return buf[:len(buf)-paddingLen], nil
}
