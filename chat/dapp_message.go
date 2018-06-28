package chat

import (
	"encoding/hex"
	"github.com/Bit-Nation/x3dh"
	"golang.org/x/crypto/ed25519"
	"time"
)

// send a message to a chat partner
func (c *Chat) CreateDAppMessage(msg string, secretID string, sec x3dh.SharedSecret, receiver ed25519.PublicKey) (Message, error) {

	// create doublerachet cipher text
	encryptedMessage, err := c.encryptMessage(sec, []byte(msg))
	if err != nil {
		return Message{}, err
	}

	m := Message{
		Type:                 "DAPP_MESSAGE",
		SendAt:               time.Now(),
		UsedSecretRef:        secretID,
		DoubleratchetMessage: encryptedMessage,
		Receiver:             hex.EncodeToString(receiver[:]),
	}

	// sign message
	err = m.Sign(c.km)
	if err != nil {
		return Message{}, err
	}

	return m, err

}
