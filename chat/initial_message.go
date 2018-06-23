package chat

import (
	"encoding/hex"
	"errors"
	"time"
	"crypto/rand"

	x3dh "github.com/Bit-Nation/x3dh"
	ed25519 "golang.org/x/crypto/ed25519"
)

var randSource = rand.Reader

// a shared secret contains
// an id + the x3dh secret
type SharedSecret struct {
	ID []byte
	Secret x3dh.SharedSecret
}

type Initialisation struct {
	Msg    Message `json:"message"`
	Secret string  `json:"secret"`
}

func (c *Chat) InitializeChat(idPubKey ed25519.PublicKey, preKeyBundle PanthalassaPreKeyBundle) (Message, x3dh.InitializedProtocol, error) {

	// init the protocol
	ip, err := c.x3dh.CalculateSecret(&preKeyBundle.PublicPart)
	if err != nil {
		return Message{}, x3dh.InitializedProtocol{}, err
	}

	// create id for shared secret
	var sid [16]byte
	if _, err := randSource.Read(sid[:]); err != nil {
		return Message{}, x3dh.InitializedProtocol{}, err
	}
	
	// create encrypted message
	msg, err := c.encryptMessage(ip.SharedSecret, sid[:])
	if err != nil {
		return Message{}, x3dh.InitializedProtocol{}, err
	}

	// construct message
	m := Message{
		Type:   "PROTOCOL_INITIALISATION",
		SendAt: time.Now(),
		AdditionalData: map[string]string{
			"used_one_time_pre_key": hex.EncodeToString(ip.UsedOneTimePreKey[:]),
			"used_signed_pre_key":   hex.EncodeToString(ip.UsedSignedPreKey[:]),
			"ephemeral_key":         hex.EncodeToString(ip.EphemeralKey[:]),
		},
		DoubleratchetMessage: msg,
		IDPubKey:             idPubKey,
	}

	// sign message
	err = m.Sign(c.km)
	if err != nil {
		return Message{}, x3dh.InitializedProtocol{}, err
	}

	return m, ip, err
}

func (c *Chat) HandleInitialMessage(m Message, keyBundlePrivate PreKeyBundlePrivate) (SharedSecret, error) {

	// message type should be: PROTOCOL_INITIALISATION
	if m.Type != "PROTOCOL_INITIALISATION" {
		return SharedSecret{}, errors.New("message must be of type PROTOCOL_INITIALISATION")
	}

	// verify signature of message
	valid, err := m.VerifySignature()
	if err != nil {
		return SharedSecret{}, err
	}
	if !valid {
		return SharedSecret{}, errors.New("invalid signature")
	}

	// get my ephemeral key
	remoteEphemeralKeyStr, exist := m.AdditionalData["ephemeral_key"]
	if !exist {
		return SharedSecret{}, errors.New("missing ephemeral_key")
	}
	remoteEphemeralKeyRaw, err := hex.DecodeString(remoteEphemeralKeyStr)
	if err != nil {
		return SharedSecret{}, err
	}
	var remoteEphemeralKey x3dh.PublicKey
	copy(remoteEphemeralKey[:], remoteEphemeralKeyRaw[:])

	// Get chat ID key
	remoteID := m.DoubleratchetMessage.Header.DH
	var remoteChatID x3dh.PublicKey
	copy(remoteChatID[:], remoteID[:])

	// calculate shared secret
	sec, err := c.x3dh.SecretFromRemote(x3dh.ProtocolInitialisation{
		RemoteIdKey:        remoteChatID,
		RemoteEphemeralKey: remoteEphemeralKey,
		MyOneTimePreKey:    &keyBundlePrivate.OneTimePreKey,
		MySignedPreKey:     keyBundlePrivate.SignedPreKey,
	})

	if err != nil {
		return SharedSecret{}, err
	}
	
	id, err := c.DecryptMessage(sec, m)
	if err != nil {
		return SharedSecret{}, err
	}
	
	return SharedSecret{
		ID: []byte(id),
		Secret: sec,
	}, nil

}
