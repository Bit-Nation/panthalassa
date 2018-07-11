package chat

import (
	"encoding/hex"
	"errors"
	"fmt"

	db "github.com/Bit-Nation/panthalassa/db"
	bpb "github.com/Bit-Nation/protobuffers"
	proto "github.com/golang/protobuf/proto"
	ed25519 "golang.org/x/crypto/ed25519"
)

// send a message
func (c *Chat) SendMessage(receiver ed25519.PublicKey, msg bpb.PlainChatMessage) error {

	var handleSendError = func(err error) error {
		updateError := c.messageDB.UpdateStatus(receiver, msg.MessageID, db.StatusFailedToSend)
		if updateError != nil {
			return errors.New(fmt.Sprintf("failed to update status with error: %s - original error: %s", updateError, err))
		}
		return err
	}

	// marshal message
	rawPlainMessage, err := proto.Marshal(&msg)
	if err != nil {
		return handleSendError(err)
	}

	// check if there is a shared secret for the receiver
	exist, err := c.sharedSecStorage.HasAny(receiver)
	if err != nil {
		return handleSendError(err)
	}

	// if we don't have a shared secret we create one
	if !exist {
		// fetch pre key bundle
		preKeyBundle, err := c.backend.FetchPreKeyBundle(receiver)
		if err != nil {
			return handleSendError(err)
		}
		// run key agreement
		initializedProtocol, err := c.x3dh.CalculateSecret(preKeyBundle)
		if err != nil {
			return handleSendError(err)
		}
		// persist shared secret
		if err := c.sharedSecStorage.Put(receiver, initializedProtocol); err != nil {
			return handleSendError(err)
		}
	}

	// fetch shared secret
	ss, err := c.sharedSecStorage.GetYoungest(receiver)
	if err != nil {
		return handleSendError(err)
	}

	hasSignedPreKey, err := c.signedPreKeyStorage.HasActive()
	if err != nil {
		return handleSendError(err)
	}
	if !hasSignedPreKey {
		signedPreKey, err := c.x3dh.NewKeyPair()
		if err != nil {
			return handleSendError(err)
		}
		if err := c.signedPreKeyStorage.Put(signedPreKey); err != nil {
			return handleSendError(err)
		}
	}

	signedPreKey, err := c.signedPreKeyStorage.GetActive()
	if err != nil {
		return handleSendError(err)
	}

	// create double ratchet session
	drSession, err := c.km.CreateDoubleRatchet(ss.X3dhSS, c.drKeyStorage, signedPreKey)
	if err != nil {
		return handleSendError(err)
	}

	// encrypt message
	drMessage := drSession.RatchetEncrypt(rawPlainMessage, nil)
	if err != nil {
		return handleSendError(err)
	}

	// fetch sender public key
	senderIdPubStr, err := c.km.IdentityPublicKey()
	if err != nil {
		return handleSendError(err)
	}
	sender, err := hex.DecodeString(senderIdPubStr)
	if err != nil {
		return handleSendError(err)
	}

	// construct chat message
	msgToSend := bpb.ChatMessage{
		MessageID: []byte(msg.MessageID),
		Receiver:  receiver,
		Message: &bpb.DoubleRatchedMsg{
			DoubleRatchetPK: drMessage.Header.DH[:],
			N:               drMessage.Header.N,
			Pn:              drMessage.Header.PN,
			CipherText:      drMessage.Ciphertext,
		},
		Sender: sender,
	}

	// attach information to message that will allow receiver
	// to derive shared secret
	if !ss.Accepted {
		msgToSend.EphemeralKey = ss.EphemeralKey[:]
		if ss.UsedOneTimePreKey != nil {
			msgToSend.OneTimePreKey = ss.UsedOneTimePreKey[:]
		}
		msgToSend.SignedPreKey = ss.UsedSignedPreKey[:]
		eks, err := c.km.IdentitySign(ss.EphemeralKey[:])
		if err != nil {
			return handleSendError(err)
		}
		msgToSend.EphemeralKeySignature = eks
		msgToSend.SharedSecretCreationDate = ss.CreatedAt.Unix()
	}

	// send message to the backend
	err = c.backend.SubmitMessage(msgToSend)
	if err != nil {
		return handleSendError(err)
	}

	return c.messageDB.UpdateStatus(receiver, msg.MessageID, db.StatusSent)
}
