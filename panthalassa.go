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

//Create an new instance of panthalassa
func NewPanthalassa(keyStore, pw string) (*Panthalassa, error) {

	km, err := keyManager.Open(keyStore, pw)
	if err != nil {
		return &Panthalassa{}, nil
	}

	return &Panthalassa{
		km: km,
	}, nil
}
