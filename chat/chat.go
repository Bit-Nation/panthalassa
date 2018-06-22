package chat

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"errors"

	aes "github.com/Bit-Nation/panthalassa/crypto/aes"
	keyManager "github.com/Bit-Nation/panthalassa/keyManager"
	x3dh "github.com/Bit-Nation/x3dh"
	doubleratchet "github.com/tiabc/doubleratchet"
)

const ProtocolName = "pangea-chat"

type Chat struct {
	doubleRachetKeyStore doubleratchet.KeysStorage
	x3dh                 x3dh.X3dh
	km                   *keyManager.KeyManager
}

type Config struct {
	// url of the websocket to connect to
	WSUrl string
	// url of the http endpoint
	HttpEndpoint string
	// access token for http endpoint and websocket
	AccessKey string
}

// create a new chat
func New(chatIdentityKey x3dh.KeyPair, km *keyManager.KeyManager, dRKeyStore doubleratchet.KeysStorage) (Chat, error) {

	c := x3dh.NewCurve25519(rand.Reader)

	x := x3dh.New(&c, sha512.New(), ProtocolName, chatIdentityKey)

	return Chat{
		doubleRachetKeyStore: dRKeyStore,
		x3dh:                 x,
		km:                   km,
	}, nil

}

// create a new pre key bundle
func (c *Chat) NewPreKeyBundle() (PanthalassaPreKeyBundle, error) {

	// @todo usually this should be kept for longer time. Need to change that later on.
	signedPreKey, err := c.x3dh.NewKeyPair()
	if err != nil {
		return PanthalassaPreKeyBundle{}, nil
	}

	oneTimePreKey, err := c.x3dh.NewKeyPair()
	if err != nil {
		return PanthalassaPreKeyBundle{}, err
	}

	chatIdKey, err := c.km.ChatIdKeyPair()
	if err != nil {
		return PanthalassaPreKeyBundle{}, err
	}

	idPubKeyStr, err := c.km.IdentityPublicKey()

	//unmarshal public key
	decodedPubIdKey, err := hex.DecodeString(idPubKeyStr)
	if err != nil {
		return PanthalassaPreKeyBundle{}, err
	}

	preKeyB := PanthalassaPreKeyBundle{
		PublicPart: PreKeyBundlePublic{
			BChatIdentityKey: chatIdKey.PublicKey,
			BSignedPreKey:    signedPreKey.PublicKey,
			BOneTimePreKey:   oneTimePreKey.PublicKey,
			BIdentityKey:     decodedPubIdKey,
		},
		PrivatePart: PreKeyBundlePrivate{
			OneTimePreKey: oneTimePreKey.PrivateKey,
			SignedPreKey:  signedPreKey.PrivateKey,
		},
	}

	if err := preKeyB.PublicPart.Sign(*c.km); err != nil {
		return PanthalassaPreKeyBundle{}, err
	}

	return preKeyB, nil

}

// export X3DH secret
func EncryptX3DHSecret(b x3dh.SharedSecret, km *keyManager.KeyManager) (aes.CipherText, error) {
	return km.AESEncrypt(b[:])
}

// decrypt x3dh secret
func DecryptX3DHSecret(secret aes.CipherText, km *keyManager.KeyManager) (x3dh.SharedSecret, error) {
	sec, err := km.AESDecrypt(secret)
	if err != nil {
		return x3dh.SharedSecret{}, err
	}
	if len(sec) != 32 {
		return x3dh.SharedSecret{}, errors.New("x3dh shared secret must have 32 bytes")
	}
	sh := [32]byte{}
	copy(sh[:], sec[:])
	return sh, nil
}
