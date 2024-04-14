package Utils

import (
	"crypto/aes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"
)

// GetTransmissionKey gets a key based on user object id. This key is used to decrypt user purchased
func GetTransmissionKey(objectID string) []byte {
	salt := "rzKKLjXkr3NpH994oYys63mSBrrmGBRComXEpXEg"
	combinedString := salt + objectID

	hash := sha256.Sum256([]byte(combinedString))

	return hash[:]
}

// RotaenoDecrypt decrypts rotaeno game data with a commonly used structures
func RotaenoDecrypt(data, key []byte) ([]byte, error) {
	iv := data[:aes.BlockSize]
	source := data[aes.BlockSize:]

	return AesDecrypt(source, key, iv)
}

// RotaenoEncrypt encrypt rotaeno game data with a commonly used structures
func RotaenoEncrypt(data, key []byte) ([]byte, error) {

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	encryptedData, err := AesEncrypt(data, key, iv)
	if err != nil {
		return nil, err
	}

	return append(iv, encryptedData...), nil
}

// RotaenoDecryptFromBase64 decrypts rotaeno game data with a commonly used structures encoded by base64
func RotaenoDecryptFromBase64(base64data string, key []byte) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(base64data)
	if err != nil {
		return nil, err
	}
	result, err := RotaenoDecrypt(data, key)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// RotaenoEncryptToBase64 encrypt rotaeno game data with a commonly used structures encoded by base64
func RotaenoEncryptToBase64(data, key []byte) (string, error) {
	encrypted, err := RotaenoEncrypt(data, key)
	encryptedBase64 := base64.StdEncoding.EncodeToString(encrypted)

	return encryptedBase64, err
}
