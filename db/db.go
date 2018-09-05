package db

import (
	"path/filepath"

	km "github.com/Bit-Nation/panthalassa/keyManager"
)

// get database path for key manager
func KMToDBPath(dir string, km *km.KeyManager) (string, error) {

	idPubKey, err := km.IdentityPublicKey()
	if err != nil {
		return "", err
	}

	return filepath.Abs(filepath.Join(dir, idPubKey+".db"))

}
