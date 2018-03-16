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
