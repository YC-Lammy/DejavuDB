package datastore

import (
	"crypto/aes"
	"io/ioutil"
	"os"

	json "github.com/goccy/go-json"

	"../settings"
)

func write_json(data map[string]interface{}) error {

	a, err := json.Marshal(data)
	if err != nil {
		return err
	}

	a, err = EncryptAES([]byte(settings.AES_key), a)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("file.json", a, 0777)
	if err != nil {
		return err
	}

	return nil
}

func read_json(filename string) (map[string]interface{}, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	data, err = DecryptAES([]byte(settings.AES_key), data)

	if err != nil {
		return nil, err
	}
	var value map[string]interface{}

	err = json.Unmarshal(data, &value)
	if err != nil {
		return nil, err
	}

	return value, nil
}

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
