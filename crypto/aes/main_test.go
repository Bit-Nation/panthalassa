package aes

import (
	"testing"
)

//Test the encrypt and decrypt function in one batch
func TestEncryptAndDecrypt(t *testing.T) {

	var secret, value string

	secret = "1111111111111111"
	value = "I am the value"

	cipherText, e := Encrypt(secret, value)

	if e != nil {
		t.Error(e)
	}

	rawValue, e := Decrypt(secret, cipherText)

	if e != nil {
		t.Error(e)
	}

	if rawValue != value {
		t.Errorf("Expected %s and %s to match", rawValue, value)
	}

}

func TestEncryptInvalidKeyLen(t *testing.T) {
	_, e := Encrypt("too_short", "")

	if e.Error() != "crypto/aes: invalid key size 9" {
		t.Error("too_short is only 9 bytes long. Valid key's are 16, 24, 28 bytes long")
	}

}

func TestDecryptInvalidKeyLen(t *testing.T) {
	_, e := Decrypt("too_short", "")

	if e.Error() != "crypto/aes: invalid key size 9" {
		t.Error("too_short is only 9 bytes long. Valid key's are 16, 24, 28 bytes long")
	}

}
