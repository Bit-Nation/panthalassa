package main

import (
	"testing"
)

//Test the encrypt and decrypt function in one batch
func TestEncryptAndDecrypt(t *testing.T) {
	
	var secret, value string
	
	secret = "1111111111111111"
	value = "I am the value"
	
	cipherText, e := encrypt( secret, value)

	if e != nil {
		t.Error(e)
	}
	
	rawValue, e := decrypt(secret, cipherText)
	
	if e != nil {
		t.Error(e)
	}
	
	if rawValue != value {
		t.Errorf("Expected %s and %s to match", rawValue, value)
	}
	
}