package cid

import (
	cid "github.com/ipfs/go-cid"
	mbase "github.com/multiformats/go-multibase"
	mh "github.com/multiformats/go-multihash"
)

//Check if CID is valid
func IsValidCid(id string) bool {

	_, err := cid.Decode(id)

	if err == nil {
		return true
	}

	return false

}

//Get the CID of a value with
//sha3 256 as a base64 string
func Sha256(value string) (string, error) {

	c, err := cid.NewPrefixV1(cid.Raw, mh.SHA3_256).Sum([]byte(value))

	if err != nil {
		return "", err
	}

	return c.StringOfBase(mbase.Base64)

}

//Get the CID of a value with
//sha3 512 as a base64 string
func Sha512(value string) (string, error) {

	c, err := cid.NewPrefixV1(cid.Raw, mh.SHA3_512).Sum([]byte(value))

	if err != nil {
		return "", err
	}

	return c.StringOfBase(mbase.Base64urlPad)

}
