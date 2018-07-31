package panthalassa

import (
	"encoding/hex"
	"fmt"

	api "github.com/Bit-Nation/panthalassa/api"
	chat "github.com/Bit-Nation/panthalassa/chat"
	dAppReg "github.com/Bit-Nation/panthalassa/dapp/registry"
	db "github.com/Bit-Nation/panthalassa/db"
	keyManager "github.com/Bit-Nation/panthalassa/keyManager"
	p2p "github.com/Bit-Nation/panthalassa/p2p"
	lp2pCrypto "github.com/libp2p/go-libp2p-crypto"
	peer "github.com/libp2p/go-libp2p-peer"
)

type Panthalassa struct {
	km       *keyManager.KeyManager
	upStream api.UpStream
	api      *api.API
	p2p      *p2p.Network
	dAppReg  *dAppReg.Registry
	chat     *chat.Chat
	msgDB    *db.BoltChatMessageStorage
}

//Stop the panthalassa instance
//this becomes interesting when we start
//to use the mesh network
func (p *Panthalassa) Stop() error {
	return p.p2p.Close()
}

//Export account with the given password
func (p *Panthalassa) Export(pw, pwConfirm string) (string, error) {

	// export
	store, err := p.km.Export(pw, pwConfirm)
	if err != nil {
		return "", err
	}

	// marshal key store
	rawStore, err := store.Marshal()
	if err != nil {
		return "", err
	}

	return string(rawStore), nil

}

// add friend to peer store
func (p *Panthalassa) AddContact(pubKey string) error {

	// decode public key
	rawPubKey, err := hex.DecodeString(pubKey)
	if err != nil {
		return err
	}

	// create lp2p public key
	lp2pPubKey, err := lp2pCrypto.UnmarshalEd25519PublicKey(rawPubKey)
	if err != nil {
		return err
	}

	// create ID from friend public key
	id, err := peer.IDFromPublicKey(lp2pPubKey)
	if err != nil {
		return err
	}

	// add public key to peer store
	err = p.p2p.Host.Peerstore().AddPubKey(id, lp2pPubKey)
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("added contact: %s", pubKey))

	return p.p2p.Host.Peerstore().Put(id, "contact", true)

}
