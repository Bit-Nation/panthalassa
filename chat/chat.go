package chat

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"errors"

	client "github.com/Bit-Nation/panthalassa/client"
	keyManager "github.com/Bit-Nation/panthalassa/keyManager"
	x3dh "github.com/Bit-Nation/x3dh"
	doubleratchet "github.com/tiabc/doubleratchet"
)

const ProtocolName = "pangea-chat"

type Chat struct {
	client               client.Client
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
func New(chatIdentityKey x3dh.KeyPair, km *keyManager.KeyManager, dRKeyStore doubleratchet.KeysStorage, client client.Client) (Chat, error) {

	c := x3dh.NewCurve25519(rand.Reader)

	x := x3dh.New(&c, sha512.New(), ProtocolName, chatIdentityKey)

	return Chat{
		client:               client,
		doubleRachetKeyStore: dRKeyStore,
		x3dh:                 x,
		km:                   km,
	}, nil

}

// create a new pre key bundle
func (c *Chat) NewPreKeyBundle() (LocalPreKeyBundle, error) {

	sigedPreKey := c.client.FetchSignedPreKey()

	oneTimePreKey, err := c.x3dh.NewKeyPair()
	if err != nil {
		return LocalPreKeyBundle{}, err
	}

	chatIdKey, err := c.km.ChatIdKeyPair()
	if err != nil {
		return LocalPreKeyBundle{}, err
	}

	idPubKeyStr, err := c.km.IdentityPublicKey()

	//unmarshal public key
	decodedPubIdKey, err := hex.DecodeString(idPubKeyStr)
	if err != nil {
		return LocalPreKeyBundle{}, err
	}

	preKeyB := LocalPreKeyBundle{
		BChatIdentityKey: chatIdKey.PublicKey,
		BSignedPreKey:    sigedPreKey.PublicKey,
		BOneTimePreKey:   oneTimePreKey.PublicKey,
		BIdentityKey:     decodedPubIdKey,
	}

	if err := preKeyB.Sign(*c.km); err != nil {
		return LocalPreKeyBundle{}, err
	}

	return preKeyB, nil

}

func (c *Chat) CreateSharedSecret(b LocalPreKeyBundle) (x3dh.InitializedProtocol, error) {
	return c.x3dh.CalculateSecret(&b)
}

// export X3DH secret
func EncryptX3DHSecret(b x3dh.SharedSecret, km *keyManager.KeyManager) (string, error) {
	secretStr := hex.EncodeToString(b[:])
	return km.AESEncrypt(secretStr)
}

// decrypt x3dh secret
func DecryptX3DHSecret(secret string, km *keyManager.KeyManager) (x3dh.SharedSecret, error) {
	sec, err := km.AESDecrypt(secret)
	if err != nil {
		return x3dh.SharedSecret{}, err
	}
	rawSecret, err := hex.DecodeString(sec)
	if len(rawSecret) != 32 {
		return x3dh.SharedSecret{}, errors.New("x3dh shared secret must have 32 bytes")
	}
	sh := [32]byte{}
	copy(sh[:], rawSecret[:])
	return sh, nil
}
