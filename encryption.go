package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"

	"golang.org/x/crypto/argon2"
)

func argonKey(key []byte) []byte {
	salt := []byte("D3G9nTEpqfAMc4mezZxdFepteNDQq8qhMs778mMWLPLbu9T7jSDsbEckUzPuxzuX")
	return argon2.IDKey(key, salt, 1, 64*1024, 4, 32)
}

func decryptAES(data []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(argonKey(key))
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
	aes, err := aes.NewCipher(argonKey(key))
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
	return ciphertext, nil
}
