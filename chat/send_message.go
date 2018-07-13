package chat

import (
	"encoding/hex"
	"errors"
	"fmt"

	prekey "github.com/Bit-Nation/panthalassa/chat/prekey"
	db "github.com/Bit-Nation/panthalassa/db"
	bpb "github.com/Bit-Nation/protobuffers"
	x3dh "github.com/Bit-Nation/x3dh"
	proto "github.com/golang/protobuf/proto"
	dr "github.com/tiabc/doubleratchet"
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

	var fetchSignedPreKey = func(userIDPubKey ed25519.PublicKey) (prekey.PreKey, error) {
		signedPreKey, err := c.userStorage.GetSignedPreKey(receiver)
		if err != nil {
			return prekey.PreKey{}, handleSendError(err)
		}

		// validate signature of signed pre key
		validSignature, err := signedPreKey.VerifySignature(userIDPubKey)
		if err != nil {
			return prekey.PreKey{}, handleSendError(err)
		}
		if !validSignature {
			return prekey.PreKey{}, handleSendError(errors.New("received invalid signature for pre key bundle"))
		}
		return signedPreKey, nil
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

	hasSignedPreKey, err := c.userStorage.HasSignedPreKey(receiver)
	if err != nil {
		return handleSendError(err)
	}

	// fetch signed pre key of chat partner if we don't have it locally
	if !hasSignedPreKey {
		err = c.refreshSignedPreKey(receiver)
		if err != nil {
			return handleSendError(err)
		}
	}

	// fetch signed pre key from storage
	signedPreKey, err := fetchSignedPreKey(receiver)
	if err != nil {
		return handleSendError(err)
	}

	// check if signed pre key expired
	expired := signedPreKey.OlderThen(SignedPreKeyValidTimeFrame)
	if expired {
		err = c.refreshSignedPreKey(receiver)
		if err != nil {
			return handleSendError(err)
		}
		// fetch signed pre key from storage
		signedPreKey, err = fetchSignedPreKey(receiver)
		if err != nil {
			return handleSendError(err)
		}
	}

	// create double ratchet session
	var drSS dr.Key
	copy(drSS[:], ss.X3dhSS[:])
	var drRK dr.Key
	copy(drRK[:], signedPreKey.PublicKey[:])

	drSession, err := dr.NewWithRemoteKey(drSS, drRK)
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
		validX3dhPub := func(pub x3dh.PublicKey) error {
			if pub == [32]byte{} {
				return errors.New("got invalid x3dh public key - seems to be empty")
			}
			if len(pub) != 32 {
				return errors.New("got invalid x3dh public key - length MUST be 32")
			}
			return nil
		}
		if ss.UsedOneTimePreKey != nil {
			if err := validX3dhPub(*ss.UsedOneTimePreKey); err != nil {
				return err
			}
			msgToSend.OneTimePreKey = ss.UsedOneTimePreKey[:]
		}
		if err := validX3dhPub(ss.UsedSignedPreKey); err != nil {
			return err
		}
		if err := validX3dhPub(ss.EphemeralKey); err != nil {
			return err
		}
		msgToSend.SignedPreKey = ss.UsedSignedPreKey[:]
		msgToSend.EphemeralKey = ss.EphemeralKey[:]
		msgToSend.EphemeralKeySignature = ss.EphemeralKeySignature
		msgToSend.SharedSecretCreationDate = ss.CreatedAt.Unix()
	}

	// send message to the backend
	err = c.backend.SubmitMessage(msgToSend)
	if err != nil {
		return handleSendError(err)
	}

	return c.messageDB.UpdateStatus(receiver, msg.MessageID, db.StatusSent)
}
