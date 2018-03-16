package panthalassa

import (
	"crypto/rand"
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
