package chat

import (
	"crypto/sha256"

	"encoding/json"
	"github.com/Bit-Nation/panthalassa/keyManager"
	x3dh "github.com/Bit-Nation/x3dh"
	ed25519 "golang.org/x/crypto/ed25519"
)

// this is a local representation of the pre key bundle
type LocalPreKeyBundle struct {
	BChatIdentityKey x3dh.PublicKey    `json:"chat_identity_key"`
	BSignedPreKey    x3dh.PublicKey    `json:"signed_pre_key"`
	BOneTimePreKey   x3dh.PublicKey    `json:"one_time_pre_key"`
	BIdentityKey     ed25519.PublicKey `json:"identity_key"`
	BSignature       []byte            `json:"identity_key_signature"`
}

// concat profile
func (b *LocalPreKeyBundle) hashBundle() []byte {

	// concat profile information
	c := append(b.BChatIdentityKey[:], b.BSignedPreKey[:]...)
	c = append(c, b.BOneTimePreKey[:]...)
	c = append(c, b.BIdentityKey...)

	return sha256.New().Sum(c)

}

func (b *LocalPreKeyBundle) IdentityKey() x3dh.PublicKey {
	return b.BChatIdentityKey
}

func (b *LocalPreKeyBundle) SignedPreKey() x3dh.PublicKey {
	return b.BSignedPreKey
}

func (b *LocalPreKeyBundle) OneTimePreKey() *x3dh.PublicKey {
	return &b.BOneTimePreKey
}

func (b *LocalPreKeyBundle) ValidSignature() bool {
	return ed25519.Verify(b.BIdentityKey, b.hashBundle(), b.BSignature)
}

// sign profile with given private key
func (b *LocalPreKeyBundle) Sign(km keyManager.KeyManager) error {

	var err error
	b.BSignature, err = km.IdentitySign(b.hashBundle())

	return err

}

// marshal the pre key bundle
func (b *LocalPreKeyBundle) Marshal() ([]byte, error) {
	return json.Marshal(b)
}

// unmarshal pre key bundle
func UnmarshalPreKeyBundle(preKeyBundle []byte) (LocalPreKeyBundle, error) {

	var b LocalPreKeyBundle
	err := json.Unmarshal(preKeyBundle, &b)
	if err != nil {
		return LocalPreKeyBundle{}, nil
	}

	return b, nil

}
