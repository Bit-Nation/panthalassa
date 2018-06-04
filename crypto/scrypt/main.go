package scrypt

import (
	"crypto/rand"
	"encoding/json"
	"errors"

	aes "github.com/Bit-Nation/panthalassa/crypto/aes"
	scrypt "gx/ipfs/QmaPHkZLbQQbvcyavn8q1GFHg6o6yeceyHFSJ3Pjf3p3TQ/go-crypto/scrypt"
)

const n = 16384
const r = 8
const p = 1
const saltLength = 50
const keyLength = 32

type Key struct {
	N      int    `json:"n"`
	R      int    `json:"r"`
	P      int    `json:"p"`
	KeyLen int    `json:"key_len"`
	Salt   []byte `json:"salt"`
	key    aes.Secret
}

type CipherText struct {
	CipherText string `json:"cipher_text"`
	ScryptKey  Key    `json:"scrypt_key"`
}

// exports CipherText as json
func (s *CipherText) Export() (string, error) {

	jsonData, err := json.Marshal(s)

	if err != nil {
		return "", err
	}

	return string(jsonData), nil

}

// derives a key out of a password
func makeScryptKey(pw string) (Key, error) {

	salt := make([]byte, saltLength)

	rand.Read(salt)

	key, err := scrypt.Key([]byte(pw), salt, n, r, p, keyLength)
	if err != nil {
		return Key{}, err
	}
	if len(key) != 32 {
		return Key{}, errors.New("key must be of length 32 in order to be used with AES")
	}

	var aesSecret aes.Secret
	copy(aesSecret[:], key[:])

	sV := Key{
		N:      n,
		R:      r,
		P:      p,
		KeyLen: keyLength,
		Salt:   salt,
		key:    aesSecret,
	}

	return sV, nil
}

//Create new ScryptCipherText
func NewCipherText(data string, password string) (string, error) {

	derivedKey, err := makeScryptKey(password)
	if err != nil {
		return "", err
	}

	cipherText, err := aes.Encrypt(data, derivedKey.key)

	cipher := CipherText{
		CipherText: cipherText,
		ScryptKey:  derivedKey,
	}

	return cipher.Export()

}

func DecryptCipherText(data, password string) (string, error) {

	var c CipherText

	if err := json.Unmarshal([]byte(data), &c); err != nil {
		return "", err
	}

	key, err := scrypt.Key([]byte(password), []byte(c.ScryptKey.Salt), c.ScryptKey.N, c.ScryptKey.R, c.ScryptKey.P, c.ScryptKey.KeyLen)
	if err != nil {
		return "", err
	}

	var AESSecret aes.Secret
	copy(AESSecret[:], key[:32])

	return aes.Decrypt(c.CipherText, AESSecret)
}
