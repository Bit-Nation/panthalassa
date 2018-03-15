package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"io"
	"crypto/rand"
	"errors"
)

// encrypt string to base64 crypto using AES
func encrypt(key string, text string) (string, error) {
	// key := []byte(keyText)
	plaintext := []byte(text)
	
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	
	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	cipherText := make([]byte, aes.BlockSize+len(plaintext))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plaintext)
	
	// convert to base64
	return base64.URLEncoding.EncodeToString(cipherText), nil
}

// decrypt from base64 to decrypted string
func decrypt(key string, cryptoText string) (string, error) {
	cipherText, _ := base64.URLEncoding.DecodeString(cryptoText)
	
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	
	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(cipherText) < aes.BlockSize {
		return "", errors.New("cipherText too short")
	}
	
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]
	
	stream := cipher.NewCFBDecrypter(block, iv)
	
	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(cipherText, cipherText)
	
	return fmt.Sprintf("%s", cipherText), nil
}