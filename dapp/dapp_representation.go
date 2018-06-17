package dapp

import (
	"bytes"
	"errors"

	mh "github.com/multiformats/go-multihash"
	ed25519 "golang.org/x/crypto/ed25519"
)

var InvalidSignature = errors.New("failed to verify signature for DApp")

// JSON Representation of published DApp
type JsonRepresentation struct {
	Name               string `json:"name"`
	Code               string `json:"code"`
	SignaturePublicKey []byte `json:"signature_public_key"`
	Signature          []byte `json:"signature"`
}

// hash the published DApp
func (r JsonRepresentation) Hash() ([]byte, error) {

	buff := bytes.NewBuffer([]byte(r.Name))

	if _, err := buff.Write([]byte(r.Code)); err != nil {
		return nil, err
	}

	if _, err := buff.Write(r.SignaturePublicKey); err != nil {
		return nil, err
	}

	multiHash, err := mh.Sum(buff.Bytes(), mh.SHA3_256, -1)
	if err != nil {
		return nil, err
	}

	return multiHash, nil

}

// verify if this published DApp
// was signed with the attached public key
func (r JsonRepresentation) VerifySignature() (bool, error) {

	hash, err := r.Hash()
	if err != nil {
		return false, err
	}

	return ed25519.Verify(r.SignaturePublicKey, hash, r.Signature), nil

}
