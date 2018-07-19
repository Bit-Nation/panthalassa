package chat

import (
	"time"

	backend "github.com/Bit-Nation/panthalassa/backend"
	db "github.com/Bit-Nation/panthalassa/db"
	keyManager "github.com/Bit-Nation/panthalassa/keyManager"
	x3dh "github.com/Bit-Nation/x3dh"
	dr "github.com/tiabc/doubleratchet"
)

const (
	SignedPreKeyValidTimeFrame = time.Hour * 24 * 60
)

type Chat struct {
	messageDB            db.ChatMessageStorage
	backend              backend.Backend
	sharedSecStorage     db.SharedSecretStorage
	x3dh                 *x3dh.X3dh
	km                   *keyManager.KeyManager
	drKeyStorage         dr.KeysStorage
	signedPreKeyStorage  db.SignedPreKeyStorage
	oneTimePreKeyStorage db.OneTimePreKeyStorage
	userStorage          db.UserStorage
}
