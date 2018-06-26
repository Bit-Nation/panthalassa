package chat

import (
	"crypto/sha256"
	"encoding/json"
	"time"

	"encoding/hex"
	"github.com/Bit-Nation/panthalassa/keyManager"
	"github.com/tiabc/doubleratchet"
	"golang.org/x/crypto/ed25519"
	"sort"
)

type Message struct {
	Type                 string                `json:"type"`
	SendAt               time.Time             `json:"timestamp"`
	UsedSecretRef        string                `json:"used_secret"`
	AdditionalData       map[string]string     `json:"additional_data"`
	DoubleratchetMessage doubleratchet.Message `json:"doubleratchet_message"`
	Signature            []byte                `json:"signature"`
	IDPubKey             string                `json:"id_public_key"`
}

// hash the message data. Exclude signature
func (m *Message) hashData() ([]byte, error) {

	h := sha256.New()

	_, err := h.Write([]byte(m.Type))
	if err != nil {
		return nil, err
	}

	_, err = h.Write([]byte(m.SendAt.String()))
	if err != nil {
		return nil, err
	}

	items := []string{}

	for k, v := range m.AdditionalData {
		items = append(items, k)
		items = append(items, v)
	}

	sort.Strings(items)

	for _, v := range items {
		h.Write([]byte(v))
	}

	_, err = h.Write(m.DoubleratchetMessage.Header.Encode())
	if err != nil {
		return nil, err
	}

	_, err = h.Write(m.DoubleratchetMessage.Ciphertext)
	if err != nil {
		return nil, err
	}

	return h.Sum(nil), nil

}

// marshal message
func (m *Message) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

// sign the message with your identity key
func (m *Message) Sign(km *keyManager.KeyManager) error {
	h, err := m.hashData()
	if err != nil {
		return err
	}
	sig, err := km.IdentitySign(h)
	m.Signature = sig
	return err
}

// verify signature of message
func (m *Message) VerifySignature() (bool, error) {
	h, err := m.hashData()
	if err != nil {
		return false, err
	}

	k, err := hex.DecodeString(m.IDPubKey)
	if err != nil {
		return false, err
	}

	return ed25519.Verify(k, h, m.Signature), nil
}
