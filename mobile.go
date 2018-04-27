package panthalassa

import "errors"

var p *panthalassa

//Create a new panthalassa instance
func NewPanthalassa(accountStore, pw string) error {

	if p != nil {
		return errors.New("you need to call Stop first in order to destroy the old instance")
	}

	instance, err := newPanthalassa(accountStore, pw)

	if err != nil {
		return err
	}

	p = instance

	return nil
}

//Stop the current panthalassa instnace
func Stop() error {

	if p == nil {
		return errors.New("you have to start panthalassa first")
	}

	if err := p.Stop(); err != nil {
		return err
	}

	p = nil

	return nil

}

//Get ethereum private key of current instance
func EthereumPrivateKey() (string, error) {

	if p == nil {
		return "", errors.New("you have to start panthalassa first")
	}

	return p.EthereumPrivateKey()

}

//Export the account storage
func Export(pw, pwConfirm string) (string, error) {

	if p == nil {
		return "", errors.New("you have to start panthalassa first")
	}

	return p.Export(pw, pwConfirm)

}
