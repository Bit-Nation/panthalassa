package chat

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"

	client "github.com/Bit-Nation/panthalassa/client"
	keyManager "github.com/Bit-Nation/panthalassa/keyManager"
	x3dh "github.com/Bit-Nation/x3dh"
	ws "github.com/gorilla/websocket"
)

const ProtocolName = "pangea-chat"

type Chat struct {
	client               *client.Client
	doubleRachetKeyStore *PangeaDoubleRachedKeyStore
	x3dh                 x3dh.X3dh
	km                   *keyManager.KeyManager
	ws                   *ws.Conn
	conf                 Config
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
func New(chatIdentityKey x3dh.KeyPair, km *keyManager.KeyManager, dRKeyStore *PangeaDoubleRachedKeyStore, client *client.Client, conf Config) (Chat, error) {

	c := x3dh.NewCurve25519(rand.Reader)

	x := x3dh.New(&c, sha512.New(), ProtocolName, chatIdentityKey)

	// create websocket connection
	wsConn, _, err := ws.DefaultDialer.Dial(conf.WSUrl, nil)
	if err != nil {
		return Chat{}, err
	}

	return Chat{
		client:               client,
		doubleRachetKeyStore: dRKeyStore,
		x3dh:                 x,
		km:                   km,
		ws:                   wsConn,
		conf:                 conf,
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

// publish a message
func (c *Chat) publishMessage(m Message) error {

	// marshal message
	rawMsg, err := m.Marshal()
	if err != nil {
		return err
	}

	// send message
	return c.ws.WriteMessage(ws.TextMessage, rawMsg)

}

func (c *Chat) CreateSharedSecret(b LocalPreKeyBundle) (x3dh.InitializedProtocol, error) {
	return c.x3dh.CalculateSecret(&b)
}
