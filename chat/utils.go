package chat

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"

	db "github.com/Bit-Nation/panthalassa/db"
	bpb "github.com/Bit-Nation/protobuffers"
	x3dh "github.com/Bit-Nation/x3dh"
	mh "github.com/multiformats/go-multihash"
	dr "github.com/tiabc/doubleratchet"
	ed25519 "golang.org/x/crypto/ed25519"
)

func (c *Chat) MarkMessagesAsRead(chatID int) error {
	return c.chatStorage.ReadMessages(chatID)
}

type drDhPair struct {
	x3dhPair x3dh.KeyPair
}

func (p *drDhPair) PrivateKey() dr.Key {
	var k dr.Key
	copy(k[:], p.x3dhPair.PrivateKey[:])
	return k
}

func (p *drDhPair) PublicKey() dr.Key {
	var k dr.Key
	copy(k[:], p.x3dhPair.PublicKey[:])
	return k
}

// update the local signed pre key for a given id
func (c *Chat) refreshSignedPreKey(idPubKey ed25519.PublicKey) error {

	// fetch signed pre key bundle
	signedPreKey, err := c.backend.FetchSignedPreKey(idPubKey)
	if err != nil {
		return err
	}

	// verify signature of signed pre key bundle
	validSig, err := signedPreKey.VerifySignature(idPubKey)
	if err != nil {
		return err
	}
	if !validSig {
		return errors.New("signed pre key signature is invalid")
	}

	// check if signed pre key didn't expire
	expired := signedPreKey.OlderThan(db.SignedPreKeyValidTimeFrame)
	if expired {
		return errors.New("signed pre key expired")
	}

	return c.userStorage.PutSignedPreKey(idPubKey, signedPreKey)

}

// generate shared secret id
func sharedSecretID(sender, receiver ed25519.PublicKey, sharedSecretID []byte) (mh.Multihash, error) {
	b := bytes.NewBuffer(sender)
	if _, err := b.Write(receiver); err != nil {
		return nil, err
	}
	if _, err := b.Write(sharedSecretID); err != nil {
		return nil, err
	}
	return mh.Sum(b.Bytes(), mh.SHA3_256, -1)
}

// create identification
func sharedSecretInitID(sender, receiver ed25519.PublicKey, msg bpb.ChatMessage) (mh.Multihash, error) {
	if len(sender) != 32 {
		return nil, errors.New("sender must be 32 bytes long")
	}
	if len(receiver) != 32 {
		return nil, errors.New("receiver must be 32 bytes long")
	}
	b := bytes.Buffer{}
	if _, err := b.Write(sender); err != nil {
		return nil, err
	}
	if _, err := b.Write(receiver); err != nil {
		return nil, err
	}
	if _, err := b.Write(msg.SenderChatIDKey); err != nil {
		return nil, err
	}
	if _, err := b.Write(msg.SignedPreKey); err != nil {
		return nil, err
	}
	if _, err := b.Write(msg.OneTimePreKey); err != nil {
		return nil, err
	}
	return mh.Sum(b.Bytes(), mh.SHA2_256, -1)
}

// hash message
func hashChatMessage(msg bpb.ChatMessage) (mh.Multihash, error) {

	b := bytes.NewBuffer(nil)

	writes := []func(b *bytes.Buffer) (int, error){
		func(b *bytes.Buffer) (int, error) {
			return b.Write(msg.OneTimePreKey)
		},
		func(b *bytes.Buffer) (int, error) {
			return b.Write(msg.SignedPreKey)
		},
		func(b *bytes.Buffer) (int, error) {
			return b.Write(msg.EphemeralKey)
		},
		func(b *bytes.Buffer) (int, error) {
			return b.Write(msg.EphemeralKeySignature)
		},
		func(b *bytes.Buffer) (int, error) {
			return b.Write(msg.SenderChatIDKey)
		},
		func(b *bytes.Buffer) (int, error) {
			return b.Write(msg.SenderChatIDKeySignature)
		},
		func(b *bytes.Buffer) (int, error) {
			return b.Write(msg.Message.DoubleRatchetPK)
		},
		func(b *bytes.Buffer) (int, error) {
			n := make([]byte, 4)
			binary.BigEndian.PutUint32(n, uint32(msg.Message.N))
			return b.Write(n)
		},
		func(b *bytes.Buffer) (int, error) {
			pn := make([]byte, 4)
			binary.BigEndian.PutUint32(pn, uint32(msg.Message.Pn))
			return b.Write(pn)
		},
		func(b *bytes.Buffer) (int, error) {
			return b.Write(msg.Message.CipherText)
		},
		func(b *bytes.Buffer) (int, error) {
			return b.Write(msg.Receiver)
		},
		func(b *bytes.Buffer) (int, error) {
			return b.Write(msg.Sender)
		},
		func(b *bytes.Buffer) (int, error) {
			return b.Write(msg.MessageID)
		},
		func(b *bytes.Buffer) (int, error) {
			return b.Write(msg.UsedSharedSecret)
		},
	}

	for _, action := range writes {
		action(b)
	}

	return mh.Sum(b.Bytes(), mh.SHA3_256, -1)

}

// turns a received plain protobuf message into a database message
func protoPlainMsgToMessage(msg *bpb.PlainChatMessage) (db.Message, error) {
	m := db.Message{
		ID:          msg.MessageID,
		Message:     msg.Message,
		CreatedAt:   msg.CreatedAt,
		Version:     uint(msg.Version),
		GroupChatID: msg.GroupChatID,
	}

	if msg.AddUserPrivChat != nil {
		m.AddUserToChat = &db.AddUserToChat{}
		m.AddUserToChat.ChatID = msg.AddUserPrivChat.ChatID
		for _, user := range msg.AddUserPrivChat.Users {
			m.AddUserToChat.Users = append(m.AddUserToChat.Users, user)
		}
		m.AddUserToChat.ChatName = msg.AddUserPrivChat.GroupName
	}

	if isDAppMessage(msg) {
		m.DApp = &db.DAppMessage{
			DAppPublicKey: msg.DAppPublicKey,
			Type:          msg.Type,
			Params:        map[string]interface{}{},
		}
		// unmarshal params
		if msg.Params != nil {
			if err := json.Unmarshal(msg.Params, &m.DApp.Params); err != nil {
				return db.Message{}, err
			}
		}
		// make sure that there is no message text
		// since the protocol doesn't allow it
		m.Message = nil
	}

	return m, nil

}

// checks if a chat message is supposed to initialize a chat
func isChatInitMessage(msg *bpb.ChatMessage) bool {
	return msg.EphemeralKey != nil && msg.SignedPreKey != nil
}

func isDAppMessage(msg *bpb.PlainChatMessage) bool {
	return len(msg.DAppPublicKey) != 0
}
