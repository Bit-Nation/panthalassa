package dapp

import (
	"crypto/sha256"
	"golang.org/x/crypto/ed25519"
)

// JSON Representation of published DApp
type JsonRepresentation struct {
	Name               string `json:"name"`
	Code               []byte `json:"code"`
	SignaturePublicKey []byte `json:"signature_public_key"`
	Signature          []byte `json:"signature"`
}

// hash the published DApp
func (r JsonRepresentation) hash() ([]byte, error) {

	h := sha256.New()

	if _, err := h.Write([]byte(r.Name)); err != nil {
		return nil, err
	}

	if _, err := h.Write([]byte(r.Code)); err != nil {
		return nil, err
	}

	if _, err := h.Write([]byte(r.SignaturePublicKey)); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

// verify if this published DApp
// was signed with the attached public key
func (r JsonRepresentation) VerifySignature() (bool, error) {

	hash, err := r.hash()
	if err != nil {
		return false, err
	}

	return ed25519.Verify(r.SignaturePublicKey, hash, r.Signature), nil

}
