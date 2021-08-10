package main

import (
	"crypto/aes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func write_json(data map[string]interface{}) error {

	a, err := json.Marshal(data)
	if err != nil {
		return err
	}

	a, err = EncryptAES([]byte(Settings.password), a)
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

	data, err = DecryptAES([]byte(Settings.password), data)

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

func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
