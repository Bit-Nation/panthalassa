package chat

import (
	"errors"
	"fmt"
	"reflect"

	db "github.com/Bit-Nation/panthalassa/db"
	bpb "github.com/Bit-Nation/protobuffers"
	x3dh "github.com/Bit-Nation/x3dh"
	proto "github.com/golang/protobuf/proto"
	dr "github.com/tiabc/doubleratchet"
	ed25519 "golang.org/x/crypto/ed25519"
)

// checks if a chat message is supposed to initialize a chat
func isChatInitMessage(msg *bpb.ChatMessage) bool {
	return msg.EphemeralKey != nil && msg.SignedPreKey != nil
}

// convert a byte sequence to an x3dh public key
func byteSliceTox3dhPub(pub []byte) (x3dh.PublicKey, error) {
	var x3dhPub x3dh.PublicKey
	if len(pub) != 32 {
		return x3dhPub, errors.New("got invalid x3dh public key (must have 32 bytes length)")
	}
	copy(x3dhPub[:], pub[:])
	return x3dhPub, nil
}

// @todo when to delete private keys
func (c *Chat) handleReceivedMessage(msg *bpb.ChatMessage) error {

	// make sure sender is a valid ed25519 public key
	sender := msg.Sender
	if len(sender) != 32 {
		return errors.New(fmt.Sprintf("got invalid sender: %s", sender))
	}

	// make sure that the message double ratchet public is legit
	if len(msg.Message.DoubleRatchetPK) != 32 {
		return errors.New("got invalid double ratchet pubcic key - must have a length of 32")
	}

	var drDHKey dr.Key
	copy(drDHKey[:], msg.Message.DoubleRatchetPK)

	drMessage := dr.Message{
		Header: dr.MessageHeader{
			DH: drDHKey,
			N:  msg.Message.N,
			PN: msg.Message.Pn,
		},
		Ciphertext: msg.Message.CipherText,
	}

	// make sure the singed pre key and the double ratchet key are the
	if msg.SignedPreKey != nil && reflect.DeepEqual(msg.Message.DoubleRatchetPK, msg.SignedPreKey) {
		return errors.New("sender chat id key must be the same as the message key when using the signed pre key")
	}

	signedPreKeyPub, err := byteSliceTox3dhPub(msg.Message.DoubleRatchetPK)
	if err != nil {
		return err
	}

	// fetched signed pre key
	signedPreKey, err := c.signedPreKeyStorage.Get(signedPreKeyPub)
	if err != nil {
		return err
	}
	if signedPreKey == nil {
		return errors.New("abort chat initialization - couldn't find signed pre key")
	}

	// construct the key pair that will be used to decrypt the chat message
	x3dhSignedPreKeyPair := x3dh.KeyPair{
		PublicKey:  signedPreKeyPub,
		PrivateKey: *signedPreKey,
	}

	// decrypt message and persist it
	var decryptMessage = func(sharedSec x3dh.SharedSecret) error {

		var drSharedSec dr.Key
		copy(drSharedSec[:], sharedSec[:])

		// double ratchet session
		drSession, err := dr.New(drSharedSec, &drDhPair{x3dhPair: x3dhSignedPreKeyPair})
		if err != nil {
			return err
		}

		// decrypt message
		rawMessage, err := drSession.RatchetDecrypt(drMessage, nil)
		if err != nil {
			return err
		}

		// unmarshal message
		plainMessage := bpb.PlainChatMessage{}
		if err := proto.Unmarshal(rawMessage, &plainMessage); err != nil {
			return err
		}

		// persist message
		return c.messageDB.PersistMessage(sender, plainMessage)

	}

	// handle chat init message
	if isChatInitMessage(msg) {

		// make sure ephemeralKey is really from sender
		if !ed25519.Verify(sender, msg.EphemeralKey, msg.EphemeralKeySignature) {
			return errors.New("aborted chat initialization - invalid ephemeral key")
		}

		// make sure chat id key is really from sender
		if !ed25519.Verify(sender, msg.SenderChatIDKey, msg.SenderChatIDKeySignature) {
			return errors.New("aborted chat initialization - invalid chat id key")
		}

		// fetch shared secret based on chat init params
		sharedSecret, err := c.sharedSecStorage.SecretForChatInitMsg(msg)
		if err != nil {
			return err
		}

		// decrypt the message like we are used to when a secret exist
		if sharedSecret != nil {
			return decryptMessage(sharedSecret.X3dhSS)
		}

		// fetched signed pre key
		signedPreKey, err := c.signedPreKeyStorage.Get(msg.SignedPreKey)
		if err != nil {
			return err
		}
		if signedPreKey == nil {
			return errors.New("abort chat initialization - couldn't find signed pre key")
		}

		// fetch used one time pre key
		oneTimePreKey, err := c.oneTimePreKeyStorage.Get(msg.OneTimePreKey)
		if err != nil {
			return err
		}

		x3dhEphemeralKey, err := byteSliceTox3dhPub(msg.EphemeralKey)
		if err != nil {
			return err
		}

		remoteIdKey, err := byteSliceTox3dhPub(msg.SenderChatIDKey)
		if err != nil {
			return err
		}

		// derive shared x3dh secret
		sharedX3dhSec, err := c.x3dh.SecretFromRemote(x3dh.ProtocolInitialisation{
			RemoteIdKey:        remoteIdKey,
			RemoteEphemeralKey: x3dhEphemeralKey,
			MyOneTimePreKey:    oneTimePreKey,
			MySignedPreKey:     *signedPreKey,
		})
		if err != nil {
			return err
		}

		// x3dh shared secret to double ratchet shared secret
		var drSharedSec dr.Key
		copy(drSharedSec[:], sharedX3dhSec[:])

		// double ratchet session
		drSession, err := dr.New(drSharedSec, &drDhPair{x3dhPair: x3dhSignedPreKeyPair}, dr.WithKeysStorage(c.drKeyStorage))
		if err != nil {
			return err
		}
		protoMessage, err := drSession.RatchetDecrypt(drMessage, nil)
		if err != nil {
			return err
		}

		// persist shared secret in accepted mode
		err = c.sharedSecStorage.Put(msg.Sender, db.SharedSecret{
			X3dhSS:   sharedX3dhSec,
			Accepted: true,
		})
		if err != nil {
			return err
		}

		// unmarshal protobuf message
		plainMsg := bpb.PlainChatMessage{}
		if err := proto.Unmarshal(protoMessage, &plainMsg); err != nil {
			return err
		}

		return c.messageDB.PersistMessage(msg.Sender, plainMsg)

	}

}
