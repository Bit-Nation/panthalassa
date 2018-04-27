package panthalassa

import (
	"github.com/Bit-Nation/panthalassa/keyManager"
)

type panthalassa struct {
	km *keyManager.KeyManager
}

//Get ethereum private key of panthalassa instance
func (p panthalassa) EthereumPrivateKey() (string, error) {
	return p.km.GetEthereumPrivateKey()
}

//Stop the panthalassa instance
//this becomes interesting when we start
//to use the mesh network
func (p *panthalassa) Stop() error {
	return nil
}

//Export account with the given password
func (p *panthalassa) Export(pw, pwConfirm string) (string, error) {
	return p.km.Export(pw, pwConfirm)
}

//Create an new instance of panthalassa
func newPanthalassa(encryptedAccount, pw string) (*panthalassa, error) {

	km, err := keyManager.OpenWithPassword(encryptedAccount, pw)
	if err != nil {
		return &panthalassa{}, nil
	}

	return &panthalassa{
		km: km,
	}, nil
}
