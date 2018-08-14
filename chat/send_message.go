package chat

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	prekey "github.com/Bit-Nation/panthalassa/chat/prekey"
	db "github.com/Bit-Nation/panthalassa/db"
	bpb "github.com/Bit-Nation/protobuffers"
	x3dh "github.com/Bit-Nation/x3dh"
	proto "github.com/golang/protobuf/proto"
	dr "github.com/tiabc/doubleratchet"
	ed25519 "golang.org/x/crypto/ed25519"
)

// send a message
func (c *Chat) SendMessage(receiver ed25519.PublicKey, dbMessage db.Message) error {

	// validate database message
	if err := db.ValidMessage(dbMessage); err != nil {
		return err
	}

	// create plain message from database message
	plainMessage := bpb.PlainChatMessage{
		CreatedAt: dbMessage.CreatedAt,
		Message:   dbMessage.Message,
		MessageID: dbMessage.ID,
		// this version is NOT the same as the version from the database message
		Version: 1,
	}

	// in the case this is a DApp message,
	// we have to add the props to the protobuf message
	if dbMessage.DApp != nil {
		plainMessage.Type = dbMessage.DApp.Type
		plainMessage.DAppPublicKey = dbMessage.DApp.DAppPublicKey
		// marshal props
		params, err := json.Marshal(dbMessage.DApp.Params)
		if err != nil {
			return err
		}
		plainMessage.Params = params
		// make sure there is no plain text in the message
		// since that's not allowed
		plainMessage.Message = nil
	}

	var handleSendError = func(err error) error {
		updateError := c.messageDB.UpdateStatus(receiver, dbMessage.DatabaseID, db.StatusFailedToSend)
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
		return *signedPreKey, nil
	}

	// @todo we should validate the plain message

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

		// ephemeral key signature
		eks, err := c.km.IdentitySign(initializedProtocol.EphemeralKey[:])
		if err != nil {
			return err
		}

		// shared secret base ID
		ssBaseID := make([]byte, 32)
		if _, err := rand.Read(ssBaseID); err != nil {
			return err
		}

		// construct shared secret
		ss := db.SharedSecret{
			ID:                    ssBaseID,
			X3dhSS:                initializedProtocol.SharedSecret,
			Accepted:              false,
			CreatedAt:             time.Now(),
			UsedOneTimePreKey:     initializedProtocol.UsedOneTimePreKey,
			UsedSignedPreKey:      initializedProtocol.UsedSignedPreKey,
			EphemeralKey:          initializedProtocol.EphemeralKey,
			EphemeralKeySignature: eks,
			BaseID:                ssBaseID,
		}

		// persist shared secret
		if err := c.sharedSecStorage.Put(receiver, ss); err != nil {
			return handleSendError(err)
		}
	}

	// fetch shared secret
	ss, err := c.sharedSecStorage.GetYoungest(receiver)
	if err != nil {
		return handleSendError(err)
	}
	if ss == nil {
		return errors.New("failed to fetch youngest secret")
	}

	hasSignedPreKey, err := c.userStorage.GetSignedPreKey(receiver)
	if err != nil {
		return handleSendError(err)
	}

	// fetch signed pre key of chat partner if we don't have it locally
	if hasSignedPreKey == nil {
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
	expired := signedPreKey.OlderThan(db.SignedPreKeyValidTimeFrame)
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

	// in the case the shared secret has not been accepted
	// we need to attach the shared secret base id
	if !ss.Accepted {
		if len(ss.BaseID) != 32 {
			return handleSendError(errors.New("base it is invalid - must have 32 bytes"))
		}
		plainMessage.SharedSecretBaseID = ss.BaseID
		plainMessage.SharedSecretCreationDate = ss.CreatedAt.Unix()
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

	// marshal message
	rawPlainMessage, err := proto.Marshal(&plainMessage)
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
		MessageID: []byte(dbMessage.ID),
		Receiver:  receiver,
		Message: &bpb.DoubleRatchetMsg{
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

		chatIDKeyPair, err := c.km.ChatIdKeyPair()
		if err != nil {
			return err
		}
		chatIDKeySignature, err := c.km.IdentitySign(chatIDKeyPair.PublicKey[:])
		if err != nil {
			return err
		}
		msgToSend.SenderChatIDKey = chatIDKeyPair.PublicKey[:]
		msgToSend.SenderChatIDKeySignature = chatIDKeySignature

		msgToSend.EphemeralKey = ss.EphemeralKey[:]
		msgToSend.EphemeralKeySignature = ss.EphemeralKeySignature
	}

	// make sure the base id is OK
	if len(ss.BaseID) != 32 {
		return errors.New("invalid base id - expected to be 32 bytes long")
	}

	// attach shared secret id to message
	msgToSend.UsedSharedSecret, err = sharedSecretID(sender, receiver, ss.BaseID)

	// send message to the backend
	err = c.backend.SubmitMessages([]*bpb.ChatMessage{&msgToSend})
	if err != nil {
		return handleSendError(err)
	}

	return c.messageDB.UpdateStatus(receiver, dbMessage.DatabaseID, db.StatusSent)
}
