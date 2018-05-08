package scrypt

import (
	"crypto/rand"
	"encoding/json"

	aes "github.com/Bit-Nation/panthalassa/crypto/aes"
	scrypt "golang.org/x/crypto/scrypt"
)

const ScryptN = 16384
const ScryptR = 8
const ScryptP = 1
const ScryptSaltLength = 50

type Key struct {
	N      int    `json:"n"`
	R      int    `json:"r"`
	P      int    `json:"p"`
	KeyLen int    `json:"key_len"`
	Salt   []byte `json:"salt"`
	key    []byte
}

type ScryptCipherText struct {
	CipherText string `json:"cipher_text"`
	ScryptKey  Key    `json:"scrypt_key"`
}

//Export's ScryptCipherText as json
func (s *ScryptCipherText) Export() (string, error) {

	jsonData, err := json.Marshal(s)

	if err != nil {
		return "", err
	}

	return string(jsonData), nil

}

//Derives a key out of a password
func Scrypt(pw string, keyLen int) (Key, error) {

	salt := make([]byte, ScryptSaltLength)

	rand.Read(salt)

	key, err := scrypt.Key([]byte(pw), salt, ScryptN, ScryptR, ScryptP, keyLen)

	if err != nil {
		return Key{}, err
	}

	sV := Key{
		N:      ScryptN,
		R:      ScryptR,
		P:      ScryptP,
		KeyLen: keyLen,
		Salt:   salt,
		key:    key,
	}

	return sV, nil
}

//Create new ScryptCipherText
func NewCipherText(data, key string) (string, error) {

	derivedKey, err := Scrypt(key, 32)

	if err != nil {
		return "", err
	}

	cipherText, err := aes.Encrypt(data, string(derivedKey.key))

	cipher := ScryptCipherText{
		CipherText: cipherText,
		ScryptKey:  derivedKey,
	}

	return cipher.Export()

}

func DecryptCipherText(data, key string) (string, error) {

	var c ScryptCipherText

	if err := json.Unmarshal([]byte(data), &c); err != nil {
		return "", err
	}

	dK, err := scrypt.Key([]byte(key), []byte(c.ScryptKey.Salt), c.ScryptKey.N, c.ScryptKey.R, c.ScryptKey.P, c.ScryptKey.KeyLen)
	if err != nil {
		return "", err
	}

	return aes.Decrypt(c.CipherText, string(dK))
}
