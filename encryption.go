package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"

	"golang.org/x/crypto/argon2"
)

func argonKey(key []byte, salt []byte) []byte {
	return argon2.IDKey(key, salt, 1, 64*1024, 4, 32)
}

func decryptAES(data []byte, key []byte) ([]byte, error) {
	saltLen := 64
	salt := data[:saltLen]
	data = data[saltLen:]
	c, err := aes.NewCipher(argonKey(key, salt))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}
	noonceSize := gcm.NonceSize()
	if noonceSize > len(data) {
		return nil, fmt.Errorf("data is too short to be encrypted")
	}
	nonce, cypher := data[:noonceSize], data[noonceSize:]
	plain, err := gcm.Open(nil, nonce, cypher, nil)
	if err != nil {
		return nil, err
	}
	return plain, nil
}

func encryptAES(data, key []byte) ([]byte, error) {
	salt := make([]byte, 64)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	aes, err := aes.NewCipher(argonKey(key, salt))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, err
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return append(salt, ciphertext...), nil
}
