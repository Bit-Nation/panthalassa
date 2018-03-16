package panthalassa

import (
	"testing"
)

func TestScryptValueKeyCreation(t *testing.T) {

	keyLength := 100

	sV, e := Scrypt("a", keyLength)

	if e != nil {
		t.Error(e)
	}

	if sV.N != ScryptN {
		t.Errorf("Expected N to be: %d got %d", ScryptN, sV.N)
	}

	if sV.P != ScryptP {
		t.Errorf("Expected P to be: %d got %d", ScryptP, sV.P)
	}

	if sV.R != ScryptR {
		t.Errorf("Expected R to be: %d got %d", ScryptR, sV.R)
	}

	if len(sV.Salt) != ScryptSaltLength {
		t.Errorf("Expect Salt to have lenght: %d got %d", ScryptSaltLength, len(sV.Salt))
	}

	if len(sV.key) != keyLength {
		t.Errorf("key lenght must match the length of: %d", keyLength)
	}

}

func TestScryptCipherText_Export(t *testing.T) {

	s := ScryptCipherText{}

	encryptedValue, e := s.Export()

	if e != nil {
		t.Error(e)
	}

	if encryptedValue != "{\"CipherText\":\"\",\"ScryptKey\":{\"N\":0,\"R\":0,\"P\":0,\"KeyLen\":0,\"Salt\":null}}" {
		t.Error("Unexpected result")
	}

}

func TestDecryptScryptCipherText(t *testing.T) {
	ethKey, err := NewScryptCipherText("password", "i_am_the_text")

	if err != nil {
		t.Error(err)
	}

	plainText, err := DecryptScryptCipherText("password", ethKey)

	if plainText != "i_am_the_text" {
		t.Errorf("Expected decrypted text to be: i_am_the_text - got: %s", plainText)
	}

}
