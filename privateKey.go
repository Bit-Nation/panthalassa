package panthalassa

import (
	"crypto/rand"
	"encoding/hex"
	lp2pCrypto "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
)

//Create new mesh private key
func NewMeshPrivKey(pw string) (string, error) {

	//Private key
	privKey, _, e := lp2pCrypto.GenerateKeyPair(lp2pCrypto.RSA, 4096)

	if e != nil {
		return "", nil
	}

	pkBytes, e := privKey.Bytes()

	return NewScryptCipherText(pw, string(pkBytes))

}

//Create new ethereum private key
func NewEthereumPrivKey(pw string) (string, error) {

	pk := make([]byte, 32)

	if _, e := rand.Read(pk); e != nil {
		return "", e
	}

	return NewScryptCipherText(pw, hex.EncodeToString(pk))
}
