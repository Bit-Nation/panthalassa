package panthalassa

import (
	rootKey "github.com/Bit-Nation/panthalassa/rootKey"
)

//Create's a new root key
func NewRootKey(pw string) (string, error) {

	rk, err := rootKey.NewRootKey()
	if err != nil {
		return "", err
	}

	exportedRootKey, err := rk.Export(pw)

	if err != nil {
		return "", err
	}

	return exportedRootKey, nil

}
