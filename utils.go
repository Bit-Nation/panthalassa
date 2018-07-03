package panthalassa

import (
	"encoding/hex"
	"strings"

	keyManager "github.com/Bit-Nation/panthalassa/keyManager"
	keyStore "github.com/Bit-Nation/panthalassa/keyStore"
	mnemonic "github.com/Bit-Nation/panthalassa/mnemonic"
	profile "github.com/Bit-Nation/panthalassa/profile"
	ethc "github.com/ethereum/go-ethereum/crypto"
	bip39 "github.com/tyler-smith/go-bip39"
)

//Creates an new set of encrypted account key's
func NewAccountKeys(pw, pwConfirm string) (string, error) {

	//Create mnemonic
	mn, err := mnemonic.New()
	if err != nil {
		return "", err
	}

	//Create KeyStore
	ks, err := keyStore.NewFromMnemonic(mn)
	if err != nil {
		return "", err
	}

	km := keyManager.CreateFromKeyStore(ks)

	// export store
	store, err := km.Export(pw, pwConfirm)
	if err != nil {
		return "", err
	}

	rawStore, err := store.Marshal()
	if err != nil {
		return "", err
	}

	return string(rawStore), nil

}

//Create new account store from mnemonic
//This can e.g. be used in case you need to recover your account
func NewAccountKeysFromMnemonic(mne, pw, pwConfirm string) (string, error) {

	//Create mnemonic
	mn, err := mnemonic.FromString(mne)
	if err != nil {
		return "", err
	}

	//Create key store from mnemonic
	ks, err := keyStore.NewFromMnemonic(mn)
	if err != nil {
		return "", err
	}

	//Create keyManager
	km := keyManager.CreateFromKeyStore(ks)

	store, err := km.Export(pw, pwConfirm)
	if err != nil {
		return "", err
	}

	rawStore, err := store.Marshal()
	if err != nil {
		return "", err
	}

	return string(rawStore), nil
}

//Check if mnemonic is valid
func IsValidMnemonic(mne string) bool {

	words := strings.Split(mne, " ")

	if len(words) != 24 {
		return false
	}

	return bip39.IsMnemonicValid(mne)
}

// sign profile
func SignProfileStandAlone(name, location, image, keyManagerStore, password string) (string, error) {

	store, err := keyManager.UnmarshalStore([]byte(keyManagerStore))
	if err != nil {
		return "", err
	}

	p, err := profile.SignWithKeyManagerStore(name, location, image, store, password)

	if err != nil {
		return "", err
	}

	_, err = p.SignaturesValid()
	if err != nil {
		return "", err
	}

	rawProfile, err := p.Marshal()
	if err != nil {
		return "", err
	}

	return string(rawProfile), nil

}

// converts an ethereum public key to address
func EthPubToAddress(pub string) (string, error) {
	pubRaw, err := hex.DecodeString(pub)
	if err != nil {
		return "", err
	}
	pubKey, err := ethc.DecompressPubkey(pubRaw)
	if err != nil {
		return "", err
	}
	addr := ethc.PubkeyToAddress(*pubKey)
	return addr.String(), nil
}
