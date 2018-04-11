package panthalassa

import (
	"github.com/Bit-Nation/panthalassa/keyManager"
	"github.com/Bit-Nation/panthalassa/keyStore"
)

//Creates an new set of encrypted account key's
func NewAccountKeys(pw, pwConfirm string) (string, error) {
	ks, err := keyStore.NewKeyStoreFactory()
	if err != nil {
		return "", err
	}
	km := keyManager.CreateFromKeyStore(ks)
	return km.Export(pw, pwConfirm)
}
