package bip32

import (
	"errors"
	"fmt"
	bip32 "github.com/tyler-smith/go-bip32"
	"strconv"
	"strings"
)

type Key = bip32.Key

//Create new MasterKey as specified in bip32
//where child key's can be derived
func NewMasterKey(seed []byte) (*bip32.Key, error) {
	return bip32.NewMasterKey(seed)
}

//Derive child key from give key.
//Path format must be "m/4H/3H"
func Derive(path string, key bip32.Key) (*bip32.Key, error) {

	pathElements := strings.Split(path, "/")

	if pathElements[0] != "m" {
		return &bip32.Key{}, errors.New(fmt.Sprintf("expected first element of path to be 'm'. Got: %s", pathElements[0]))
	}

	var derivedKey *bip32.Key
	for _, pathElement := range pathElements[1:] {
		if false == strings.HasSuffix(pathElement, "H") {
			return &bip32.Key{}, errors.New(fmt.Sprintf("since we only expect hardened key's there should be an H in the path element. got: %s", pathElement))
		}

		p, err := strconv.Atoi(strings.TrimSuffix(pathElement, "H"))
		if err != nil {
			return &bip32.Key{}, err
		}

		key, err := key.NewChildKey(uint32(2147483648 + p))
		if err != nil {
			return &bip32.Key{}, err
		}
		derivedKey = key
	}

	return derivedKey, nil

}
