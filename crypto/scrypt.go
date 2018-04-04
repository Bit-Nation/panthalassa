package crypto

import (
	"crypto/rand"
	"encoding/json"

	"golang.org/x/crypto/scrypt"
)

const ScryptN = 16384
const ScryptR = 8
const ScryptP = 1
const ScryptSaltLength = 50

type ScryptKey struct {
	N      int
	R      int
	P      int
	KeyLen int
	Salt   []byte
	key    []byte
}

type ScryptCipherText struct {
	CipherText string
	ScryptKey  ScryptKey
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
func Scrypt(pw string, keyLen int) (ScryptKey, error) {

	salt := make([]byte, ScryptSaltLength)

	rand.Read(salt)

	key, err := scrypt.Key([]byte(pw), salt, ScryptN, ScryptR, ScryptP, keyLen)

	if err != nil {
		return ScryptKey{}, err
	}

	sV := ScryptKey{
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
func NewScryptCipherText(pw string, data string) (string, error) {

	derivedKey, err := Scrypt(pw, 32)

	if err != nil {
		return "", err
	}

	cipherText, err := encrypt(string(derivedKey.key), data)

	cipher := ScryptCipherText{
		CipherText: cipherText,
		ScryptKey:  derivedKey,
	}

	return cipher.Export()

}

func DecryptScryptCipherText(pw string, data string) (string, error) {

	var c ScryptCipherText

	if err := json.Unmarshal([]byte(data), &c); err != nil {
		return "", err
	}

	dK, err := scrypt.Key([]byte(pw), []byte(c.ScryptKey.Salt), c.ScryptKey.N, c.ScryptKey.R, c.ScryptKey.P, c.ScryptKey.KeyLen)

	if err != nil {
		return "", err
	}

	return decrypt(string(dK), c.CipherText)
}
