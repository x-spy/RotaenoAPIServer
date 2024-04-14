package Utils

import (
	"crypto/aes"
	"crypto/cipher"
)

func AesDecrypt(source, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	decrypter := cipher.NewCBCDecrypter(block, iv)

	decrypted := make([]byte, len(source))
	copy(decrypted, source)

	decrypter.CryptBlocks(decrypted, decrypted)

	return decrypted, nil
}

func AesEncrypt(source, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	encrypter := cipher.NewCBCEncrypter(block, iv)

	encrypted := make([]byte, len(source))
	copy(encrypted, source)

	encrypter.CryptBlocks(encrypted, encrypted)

	return encrypted, nil
}
