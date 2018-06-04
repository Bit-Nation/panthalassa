package aes

//Taken from this example: https://gist.github.com/cannium/c167a19030f2a3c6adbb5a5174bea3ff

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
)

type aesCipherText struct {
	Iv         string `json:"iv"`
	CipherText string `json:"cipher_text"`
}

func (a aesCipherText) Marshal() ([]byte, error) {
	return json.Marshal(a)
}

type Secret [32]byte

// encrypt string to base64 crypto using AES
func Encrypt(plainText string, key Secret) (string, error) {

	//Create cipher block
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}

	// Never use more than 2^32 random nonces with a given key because of the risk of a repeat.
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	aesGcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	cipherText := aesGcm.Seal(nil, nonce, []byte(plainText), nil)

	ct := aesCipherText{
		Iv:         hex.EncodeToString(nonce),
		CipherText: hex.EncodeToString(cipherText),
	}

	marshaled, err := ct.Marshal()
	if err != nil {
		return "", err
	}

	return string(marshaled), nil
}

func Decrypt(cipherText string, secret Secret) (string, error) {

	//Unmarshal cipher text
	var ct aesCipherText
	if err := json.Unmarshal([]byte(cipherText), &ct); err != nil {
		return "", err
	}

	//Decode IV
	nonce, err := hex.DecodeString(ct.Iv)
	if err != nil {
		return "", err
	}

	//Decode plain cipher text
	encryptedCipherText, err := hex.DecodeString(ct.CipherText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(secret[:])
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plainText, err := aesgcm.Open(nil, nonce, encryptedCipherText, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}
