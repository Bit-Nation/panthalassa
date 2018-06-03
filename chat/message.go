package chat

import (
	"encoding/json"
	"time"

	"crypto/sha256"
	"github.com/Bit-Nation/panthalassa/keyManager"
	"github.com/tiabc/doubleratchet"
	"golang.org/x/crypto/ed25519"
)

type Message struct {
	Type                 string                `json:"type"`
	SendAt               time.Time             `json:"timestamp"`
	AdditionalData       map[string]string     `json:"additional_data"`
	DoubleratchetMessage doubleratchet.Message `json:"doubleratchet_message"`
	Signature            []byte                `json:"signature"`
	IDPubKey             ed25519.PublicKey     `json:"id_public_key"`
}

// hash the message data. Exclude signature
func (m *Message) hashData() []byte {

	b := []byte(m.Type)
	b = append(b, []byte(m.SendAt.String())...)

	for k, v := range m.AdditionalData {
		b = append(b, []byte(k)...)
		b = append(b, []byte(v)...)
	}

	b = append(b, m.DoubleratchetMessage.Header.Encode()...)
	b = append(b, m.DoubleratchetMessage.Ciphertext...)

	return sha256.New().Sum(b)

}

// marshal message
func (m *Message) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

// sign the message with your identity key
func (m *Message) Sign(km *keyManager.KeyManager) error {
	sig, err := km.IdentitySign(m.hashData())
	m.Signature = sig
	return err
}

// verify signature of message
func (m *Message) VerifySignature() bool {
	return ed25519.Verify(m.IDPubKey, m.hashData(), m.Signature)
}
