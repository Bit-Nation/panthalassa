package chat

import (
	"crypto/sha256"
	"encoding/json"

	aes "github.com/Bit-Nation/panthalassa/crypto/aes"
	keyManager "github.com/Bit-Nation/panthalassa/keyManager"
	x3dh "github.com/Bit-Nation/x3dh"
	ed25519 "golang.org/x/crypto/ed25519"
)

type PreKeyBundlePublic struct {
	BChatIdentityKey x3dh.PublicKey    `json:"chat_identity_key"`
	BSignedPreKey    x3dh.PublicKey    `json:"signed_pre_key"`
	BOneTimePreKey   x3dh.PublicKey    `json:"one_time_pre_key"`
	BIdentityKey     ed25519.PublicKey `json:"identity_key"`
	BSignature       []byte            `json:"identity_key_signature"`
}

type PreKeyBundlePrivate struct {
	OneTimePreKey x3dh.PrivateKey `json:"one_time_pre_key"`
	SignedPreKey  x3dh.PrivateKey `json:"signed_one_time_pre_key"`
}

type ExportedPreKeyBundle struct {
	PublicPart  PreKeyBundlePublic `json:"public_part"`
	PrivatePart aes.CipherText     `json:"private_part"`
}

// marshal private part of pre key bundle
func (b *PreKeyBundlePrivate) Marshal() ([]byte, error) {
	return json.Marshal(b)
}

// this is a local representation of the pre key bundle
type PanthalassaPreKeyBundle struct {
	PublicPart  PreKeyBundlePublic  `json:"public_part"`
	PrivatePart PreKeyBundlePrivate `json:"private_part"`
}

// concat profile
func (b *PreKeyBundlePublic) hashBundle() []byte {

	// concat profile information
	c := append(b.BChatIdentityKey[:], b.BSignedPreKey[:]...)
	c = append(c, b.BOneTimePreKey[:]...)
	c = append(c, b.BIdentityKey...)

	return sha256.New().Sum(c)

}

func (b *PreKeyBundlePublic) Marshal() ([]byte, error) {
	return json.Marshal(b)
}

func (b *PreKeyBundlePublic) IdentityKey() x3dh.PublicKey {
	return b.BChatIdentityKey
}

func (b *PreKeyBundlePublic) SignedPreKey() x3dh.PublicKey {
	return b.BSignedPreKey
}

func (b *PreKeyBundlePublic) OneTimePreKey() *x3dh.PublicKey {
	return &b.BOneTimePreKey
}

func (b *PreKeyBundlePublic) ValidSignature() bool {
	return ed25519.Verify(b.BIdentityKey, b.hashBundle(), b.BSignature)
}

// sign profile with given private key
func (b *PreKeyBundlePublic) Sign(km keyManager.KeyManager) error {

	var err error
	b.BSignature, err = km.IdentitySign(b.hashBundle())

	return err

}

// marshal the pre key bundle
func (b *PanthalassaPreKeyBundle) Marshal() ([]byte, error) {
	return json.Marshal(b)
}

// unmarshal pre key bundle
func UnmarshalPreKeyBundle(preKeyBundle []byte) (PanthalassaPreKeyBundle, error) {
	var b PanthalassaPreKeyBundle
	if err := json.Unmarshal(preKeyBundle, &b); err != nil {
		return PanthalassaPreKeyBundle{}, err
	}
	return b, nil
}
