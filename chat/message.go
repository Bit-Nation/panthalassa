package chat

import (
	"time"
	"encoding/json"
	"encoding/hex"
	
	"github.com/tiabc/doubleratchet"
	"github.com/Bit-Nation/panthalassa/keyManager"
	"crypto/sha256"
)

type Message struct {
	Type                 string                `json:"type"`
	SendAt               time.Time             `json:"timestamp"`
	AdditionalData       map[string]string     `json:"additional_data"`
	DoubleratchetMessage doubleratchet.Message `json:"doubleratchet_message"`
	Signature            string					`json:"signature"`
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

func (m *Message) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

// sign the message with your identity key
func (m *Message) Sign(km *keyManager.KeyManager) error {

	sig, err := km.IdentitySign(m.hashData())
	m.Signature = hex.EncodeToString(sig)
	return err
	
}
