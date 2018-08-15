package chat

import (
	"encoding/hex"
	"errors"
	"time"

	db "github.com/Bit-Nation/panthalassa/db"
	bpb "github.com/Bit-Nation/protobuffers"
	x3dh "github.com/Bit-Nation/x3dh"
	proto "github.com/golang/protobuf/proto"
	dr "github.com/tiabc/doubleratchet"
	ed25519 "golang.org/x/crypto/ed25519"
)

// convert a byte sequence to an x3dh public key
func byteSliceTox3dhPub(pub []byte) (x3dh.PublicKey, error) {
	var x3dhPub x3dh.PublicKey
	if len(pub) != 32 {
		return x3dhPub, errors.New("got invalid x3dh public key (must have 32 bytes length)")
	}
	copy(x3dhPub[:], pub[:])
	return x3dhPub, nil
}

func (c *Chat) decryptMessage(msg dr.Message, sharedSecret x3dh.SharedSecret, signedPreKey *x3dh.KeyPair) (bpb.PlainChatMessage, error) {

	// convert x3dh shared secret to double ratchet shared secret
	x3dhToDr := func(secret x3dh.SharedSecret) dr.Key {
		var k dr.Key
		copy(k[:], secret[:])
		return k
	}

	// used to decrypt the message with the given key pair
	decrypt := func(signedPreKey x3dh.KeyPair) (bpb.PlainChatMessage, error) {
		session, err := dr.New(x3dhToDr(sharedSecret), &drDhPair{x3dhPair: signedPreKey})
		if err != nil {
			return bpb.PlainChatMessage{}, err
		}
		rawMsg, err := session.RatchetDecrypt(msg, nil)
		if err != nil {
			return bpb.PlainChatMessage{}, err
		}
		protoMsg := bpb.PlainChatMessage{}
		if err := proto.Unmarshal(rawMsg, &protoMsg); err != nil {
			return bpb.PlainChatMessage{}, err
		}
		return protoMsg, nil
	}

	// in the case a signed pre key was passed in
	// take it and decrypt
	if signedPreKey != nil {
		return decrypt(*signedPreKey)
	}

	// in the case that no signed pre key was passed in
	// we fetch all our pre keys and try to decrypt.
	// trying all of them should be fine since we will only have
	// around 6 - 8 signed pre keys
	for _, signedPreKey := range c.signedPreKeyStorage.All() {
		plainChatMessage, err := decrypt(*signedPreKey)
		if err != nil {
			return bpb.PlainChatMessage{}, err
		}
		return plainChatMessage, nil
	}

	return bpb.PlainChatMessage{}, errors.New("failed to decrypt message")

}

func (c *Chat) handleReceivedMessage(msg *bpb.ChatMessage) error {

	// @todo HERE would message authentication happen if we decide to implement it
	logger.Debugf("handle received message: %s", msg)

	// make sure we don't handle our own messages
	ourIDKey, err := c.km.IdentityPublicKey()
	if err != nil {
		return err
	}
	if ourIDKey == hex.EncodeToString(msg.Sender) {
		return errors.New("in can't handle messages I created my self - this is non sense")
	}

	// make sure sender is a valid ed25519 public key
	sender := msg.Sender
	if len(sender) != 32 {
		return errors.New("sender public key too short")
	}

	// make sure that the message double ratchet public is legit
	if len(msg.Message.DoubleRatchetPK) != 32 {
		return errors.New("got invalid double ratchet public key - must have a length of 32")
	}

	// decode our id key
	ourRawIDPubKey, err := c.km.IdentityPublicKey()
	if err != nil {
		return err
	}
	ourIDPubKey, err := hex.DecodeString(ourRawIDPubKey)
	if err != nil {
		return err
	}

	drMessage := dr.Message{
		Header: dr.MessageHeader{
			DH: func() dr.Key {
				var drDHKey dr.Key
				copy(drDHKey[:], msg.Message.DoubleRatchetPK)
				return drDHKey
			}(),
			N:  msg.Message.N,
			PN: msg.Message.Pn,
		},
		Ciphertext: msg.Message.CipherText,
	}

	logger.Debugf("double ratchet message %s", drMessage)

	// handle chat init message
	if isChatInitMessage(msg) {

		logger.Debugf("message is a chat installation message")

		// make sure ephemeralKey is really from sender
		if !ed25519.Verify(sender, msg.EphemeralKey, msg.EphemeralKeySignature) {
			return errors.New("aborted chat initialization - invalid ephemeral key")
		}

		// make sure chat id key is really from sender
		if !ed25519.Verify(sender, msg.SenderChatIDKey, msg.SenderChatIDKeySignature) {
			return errors.New("aborted chat initialization - invalid chat id key")
		}

		// make sure signed pre key is valid
		if len(msg.SignedPreKey) != 32 {
			return errors.New("aborted chat initialization - invalid signed pre key")
		}

		var signedPreKey x3dh.KeyPair
		copy(signedPreKey.PublicKey[:], msg.SignedPreKey)

		logger.Debugf("used signed pre key: %x", msg.SignedPreKey)

		// get private signed pre key
		signedPreKeyPriv, err := c.signedPreKeyStorage.Get(signedPreKey.PublicKey)
		if err != nil {
			return err
		}
		signedPreKey.PrivateKey = *signedPreKeyPriv

		idInitParams, err := sharedSecretInitID(msg.Sender, ourIDPubKey, *msg)
		if err != nil {
			return err
		}

		logger.Debugf("the id based on init params is: %x", idInitParams)

		// fetch shared secret based on chat init params
		sharedSecret, err := c.sharedSecStorage.SecretForChatInitMsg(msg.Sender, idInitParams)
		if err != nil {
			return err
		}

		// decrypt the message
		if sharedSecret != nil {

			logger.Debug("found shared secret for chat init params - decrypting")

			decryptedMsg, err := c.decryptMessage(drMessage, sharedSecret.X3dhSS, &signedPreKey)
			if err != nil {
				return err
			}
			// convert proto plain message to decrypted message
			dbMessage, err := protoPlainMsgToMessage(&decryptedMsg)
			dbMessage.Sender = sender
			if err != nil {
				return err
			}
			return c.messageDB.PersistReceivedMessage(msg.Sender, dbMessage)
		}

		// fetch used one time pre key
		oneTimePreKeyPriv, err := c.oneTimePreKeyStorage.Cut(msg.OneTimePreKey)
		if err != nil {
			return err
		}

		x3dhEphemeralKey, err := byteSliceTox3dhPub(msg.EphemeralKey)
		if err != nil {
			return err
		}
		if x3dhEphemeralKey == [32]byte{} {
			return errors.New("invalid sender ephemeral id key")
		}

		remoteChatIdKey, err := byteSliceTox3dhPub(msg.SenderChatIDKey)
		if err != nil {
			return err
		}
		if remoteChatIdKey == [32]byte{} {
			return errors.New("invalid sender chat id key")
		}

		// derive shared x3dh secret
		sharedX3dhSec, err := c.x3dh.SecretFromRemote(x3dh.ProtocolInitialisation{
			RemoteIdKey:        remoteChatIdKey,
			RemoteEphemeralKey: x3dhEphemeralKey,
			MyOneTimePreKey:    oneTimePreKeyPriv,
			MySignedPreKey:     signedPreKey.PrivateKey,
		})
		if err != nil {
			return err
		}

		// x3dh shared secret to double ratchet shared secret
		var drSharedSec dr.Key
		copy(drSharedSec[:], sharedX3dhSec[:])

		// double ratchet session
		drSession, err := dr.New(drSharedSec, &drDhPair{x3dhPair: signedPreKey}, dr.WithKeysStorage(c.drKeyStorage))
		if err != nil {
			return err
		}

		// decrypt message
		logger.Debug("try to decrypt message with newly created shared secret")
		protoMessage, err := drSession.RatchetDecrypt(drMessage, nil)
		if err != nil {
			return err
		}

		// unmarshal protobuf message
		plainMsg := bpb.PlainChatMessage{}
		if err := proto.Unmarshal(protoMessage, &plainMsg); err != nil {
			return err
		}

		// make sure used share secret id exist
		if len(plainMsg.SharedSecretBaseID) != 32 {
			return errors.New("invalid used shared secret id")
		}

		// generate shared secret id
		ssID, err := sharedSecretID(sender, ourIDPubKey, plainMsg.SharedSecretBaseID)
		if err != nil {
			return err
		}

		// make sure creation date is valid
		if plainMsg.SharedSecretCreationDate == 0 {
			return errors.New("abort chat initialization - invalid shared secret creation date")
		}

		// persist shared secret in accepted mode
		err = c.sharedSecStorage.Put(sender, db.SharedSecret{
			X3dhSS: sharedX3dhSec,
			// safe as accepted since the sender initialized the chat
			Accepted:     true,
			CreatedAt:    time.Unix(plainMsg.SharedSecretCreationDate, 0),
			ID:           ssID,
			IDInitParams: idInitParams,
			BaseID:       plainMsg.SharedSecretBaseID,
		})
		if err != nil {
			return err
		}

		// convert plain protobuf message to database message
		dbMessage, err := protoPlainMsgToMessage(&plainMsg)
		dbMessage.Sender = sender
		if err != nil {
			return err
		}

		return c.messageDB.PersistReceivedMessage(sender, dbMessage)

	}

	// in the case the message don't contain a shared secret id
	// we should exit here since the sender doesn't follow our protocol
	if len(msg.UsedSharedSecret) == 0 {
		return errors.New("message is not a chat initialisation message but don't contain information about which shared secret has been used")
	}

	// fetch shared secret
	sharedSec, err := c.sharedSecStorage.Get(sender, msg.UsedSharedSecret)
	if err != nil {
		// @todo publish status that we failed to decrypt this message
		// @todo it's something that is not supposed to happen
		return err
	}

	// exit when shared secret is not found - sender doesn't follow protocol
	if sharedSec == nil {
		return errors.New("no shared secret found but is not a chat init message - sender doesn't follow protocol")
	}

	// decrypt message with fetched x3dh shared secret
	plainMsg, err := c.decryptMessage(drMessage, sharedSec.X3dhSS, nil)
	if err != nil {
		return err
	}

	// convert proto message to database message
	dbMessage, err := protoPlainMsgToMessage(&plainMsg)
	dbMessage.Sender = sender
	if err != nil {
		return err
	}

	// persist message
	if err := c.messageDB.PersistReceivedMessage(msg.Sender, dbMessage); err != nil {
		return err
	}

	// if the decryption didn't fail we want to mark
	// the shared secret as accepted
	if !sharedSec.Accepted {
		return c.sharedSecStorage.Accept(msg.Sender, sharedSec)
	}

	return nil

}
