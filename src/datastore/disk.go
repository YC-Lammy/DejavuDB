package datastore

import (
	"crypto/aes"
)

func EncryptAES(key []byte, plaintext []byte) ([]byte, error) {
	// create cipher
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// allocate space for ciphered data
	out := make([]byte, len(plaintext))

	// encrypt
	c.Encrypt(out, plaintext)
	// return hex string
	return out, nil
}

func DecryptAES(key []byte, ciphertext []byte) ([]byte, error) {

	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	pt := make([]byte, len(ciphertext))
	c.Decrypt(pt, ciphertext)

	return pt, nil
}
