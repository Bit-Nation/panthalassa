package panthalassa

import (
	"github.com/Bit-Nation/panthalassa/keyManager"
)

type Panthalassa struct {
	km *keyManager.KeyManager
}

//Get ethereum private key of panthalassa instance
func (p Panthalassa) EthereumPrivateKey() (string, error) {
	return p.km.GetEthereumPrivateKey()
}

//Stop the panthalassa instance
//this becomes interesting when we start
//to use the mesh network
func (p *Panthalassa) Stop() error {
	return nil
}

//Export account with the given password
func (p *Panthalassa) Export(pw, pwConfirm string) (string, error) {
	return p.km.Export(pw, pwConfirm)
}

//Create an new instance of panthalassa
func NewPanthalassa(encryptedAccount, pw string) (*Panthalassa, error) {

	km, err := keyManager.OpenWithPassword(encryptedAccount, pw)
	if err != nil {
		return &Panthalassa{}, nil
	}

	return &Panthalassa{
		km: km,
	}, nil
}
