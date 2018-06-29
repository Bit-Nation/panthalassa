package chat

import (
	"encoding/hex"
	"golang.org/x/crypto/ed25519"
	"time"

	x3dh "github.com/Bit-Nation/x3dh"
)

type Info struct {
	SharedSecret x3dh.SharedSecret
}

// send a message to a chat partner
func (c *Chat) CreateHumanMessage(msg string, secretID string, sec x3dh.SharedSecret, receiver ed25519.PublicKey) (Message, error) {

	// create doublerachet cipher text
	encryptedMessage, err := c.encryptMessage(sec, []byte(msg))
	if err != nil {
		return Message{}, err
	}

	myId, err := c.km.IdentityPublicKey()
	if err != nil {
		return Message{}, err
	}

	m := Message{
		Type:                 "HUMAN_MESSAGE",
		SendAt:               time.Now(),
		UsedSecretRef:        secretID,
		DoubleratchetMessage: encryptedMessage,
		Receiver:             hex.EncodeToString(receiver[:]),
		IDPubKey:             myId,
	}

	// sign message
	err = m.Sign(c.km)
	if err != nil {
		return Message{}, err
	}

	return m, err

}
