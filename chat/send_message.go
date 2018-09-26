package chat

import (
	"encoding/hex"
	"encoding/json"
	"errors"
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
		Version:     1,
		GroupChatID: dbMessage.GroupChatID,
	}

	// attach add user data
	if dbMessage.AddUserToChat != nil {
		addUserMsg := dbMessage.AddUserToChat
		groupChat, err := c.chatStorage.GetGroupChatByRemoteID(dbMessage.GroupChatID)
		if err != nil {
			return err
		}
		plainMessage.AddUserPrivChat = &bpb.PlainChatMessage_AddUserPrivGroupChat{
			Users: func() [][]byte {
				users := [][]byte{}
				for _, u := range dbMessage.AddUserToChat.Users {
					users = append(users, u)
				}
				return users
			}(),
			ChatID:    addUserMsg.ChatID,
			GroupName: groupChat.Name,
		}
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

	var fetchSignedPreKey = func(userIDPubKey ed25519.PublicKey) (*prekey.PreKey, error) {
		signedPreKey, err := c.userStorage.GetSignedPreKey(receiver)
		if err != nil {
			return nil, err
		}

		// validate signature of signed pre key
		validSignature, err := signedPreKey.VerifySignature(userIDPubKey)
		if err != nil {
			return nil, err
		}
		if !validSignature {
			return nil, errors.New("received invalid signature for signed pre key")
		}
		return signedPreKey, nil
	}

	// fetch sender public key
	senderIdPubStr, err := c.km.IdentityPublicKey()
	if err != nil {
		return err
	}
	sender, err := hex.DecodeString(senderIdPubStr)
	if err != nil {
		return err
	}

	// @todo we should validate the plain message

	// check if there is a shared secret for the receiver
	exist, err := c.sharedSecStorage.HasAny(receiver)
	if err != nil {
		return err
	}

	// if we don't have a shared secret we create one
	if !exist {

		// fetch pre key bundle
		preKeyBundle, err := c.backend.FetchPreKeyBundle(receiver)
		if err != nil {
			return err
		}

		// run key agreement
		initializedProtocol, err := c.x3dh.CalculateSecret(preKeyBundle)
		if err != nil {
			return err
		}

		// ephemeral key signature
		eks, err := c.km.IdentitySign(initializedProtocol.EphemeralKey[:])
		if err != nil {
			return err
		}

		// shared secret base ID
		id := make([]byte, 32)
		if _, err := rand.Read(id); err != nil {
			return err
		}

		// construct shared secret
		ss := db.SharedSecret{
			ID:                    id,
			Accepted:              false,
			CreatedAt:             time.Now(),
			UsedOneTimePreKey:     initializedProtocol.UsedOneTimePreKey,
			UsedSignedPreKey:      initializedProtocol.UsedSignedPreKey,
			EphemeralKey:          initializedProtocol.EphemeralKey,
			EphemeralKeySignature: eks,
			Partner:               receiver,
		}
		ss.SetX3dhSecret(initializedProtocol.SharedSecret)

		// persist shared secret
		if err := c.sharedSecStorage.Put(ss); err != nil {
			return err
		}
	}

	// fetch shared secret
	ss, err := c.sharedSecStorage.GetYoungest(receiver)
	if err != nil {
		return err
	}
	if ss == nil {
		return errors.New("failed to fetch youngest secret")
	}

	signedPreKey, err := c.userStorage.GetSignedPreKey(receiver)
	if err != nil {
		return err
	}

	// fetch signed pre key of chat partner if we don't have it locally
	if signedPreKey == nil {
		err = c.refreshSignedPreKey(receiver)
		if err != nil {
			return err
		}
	}

	// fetch signed pre key from storage
	signedPreKey, err = fetchSignedPreKey(receiver)
	if err != nil {
		return err
	}

	// check if signed pre key expired
	expired := signedPreKey.OlderThan(db.SignedPreKeyValidTimeFrame)
	if expired {
		err = c.refreshSignedPreKey(receiver)
		if err != nil {
			return err
		}
		// fetch signed pre key from storage
		signedPreKey, err = fetchSignedPreKey(receiver)
		if err != nil {
			return err
		}
	}

	// in the case the shared secret has not been accepted
	// we need to attach the shared secret base id
	if !ss.Accepted {
		if len(ss.ID) != 32 {
			return errors.New("base it is invalid - must have 32 bytes")
		}
		plainMessage.SharedSecretBaseID = ss.ID
		plainMessage.SharedSecretCreationDate = ss.CreatedAt.Unix()
	}

	// create double ratchet session
	var drSS dr.Key
	x3dhSS := ss.GetX3dhSecret()
	copy(drSS[:], x3dhSS[:])
	var drRK dr.Key
	copy(drRK[:], signedPreKey.PublicKey[:])

	drSession, err := dr.NewWithRemoteKey(drSS, drRK)
	if err != nil {
		return err
	}

	// marshal message
	rawPlainMessage, err := proto.Marshal(&plainMessage)
	if err != nil {
		return err
	}

	// encrypt message
	drMessage := drSession.RatchetEncrypt(rawPlainMessage, nil)
	if err != nil {
		return err
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
	if len(ss.ID) != 32 {
		return errors.New("invalid base id - expected to be 32 bytes long")
	}

	// attach shared secret id to message
	msgToSend.UsedSharedSecret = ss.ID

	// send message to the backend
	err = c.backend.SubmitMessages([]*bpb.ChatMessage{&msgToSend})
	if err != nil {
		return err
	}

	return nil
	// return c.messageDB.UpdateStatus(receiver, dbMessage.DatabaseID, db.StatusSent)
}
