package keyManager

import (
	"encoding/hex"
	"encoding/json"
	"errors"

	aes "github.com/Bit-Nation/panthalassa/crypto/aes"
	scrypt "github.com/Bit-Nation/panthalassa/crypto/scrypt"
	ks "github.com/Bit-Nation/panthalassa/keyStore"
	chatMigration "github.com/Bit-Nation/panthalassa/keyStore/migration/chat"
	encryptionKeyMigration "github.com/Bit-Nation/panthalassa/keyStore/migration/encryption_key"
	ethereumMigration "github.com/Bit-Nation/panthalassa/keyStore/migration/ethereum"
	identity "github.com/Bit-Nation/panthalassa/keyStore/migration/identity"
	"github.com/Bit-Nation/panthalassa/mnemonic"
	x3dh "github.com/Bit-Nation/x3dh"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	lp2pCrypto "github.com/libp2p/go-libp2p-crypto"
)

type KeyManager struct {
	keyStore ks.Store
	account  Store
}

type Store struct {
	// the password is encrypted with the mnemonic
	Password          string `json:"password"`
	EncryptedKeyStore string `json:"encrypted_key_store"`
	Version           uint8  `json:"version"`
}

//Open encrypted keystore with password
func OpenWithPassword(encryptedStore, pw string) (*KeyManager, error) {

	//unmarshal encrypted account
	var store Store
	if err := json.Unmarshal([]byte(encryptedStore), &store); err != nil {
		return &KeyManager{}, err
	}

	//Decrypt key store
	jsonKeyStore, err := scrypt.DecryptCipherText(store.EncryptedKeyStore, pw)
	if err != nil {
		return &KeyManager{}, err
	}

	//unmarshal key store
	keyStore, err := ks.UnmarshalStore(jsonKeyStore)
	if err != nil {
		return &KeyManager{}, err
	}

	return &KeyManager{
		keyStore: keyStore,
		account:  store,
	}, nil

}

//Open account with mnemonic.
//This should only be used as a backup
func OpenWithMnemonic(encryptedAccount, mnemonic string) (*KeyManager, error) {

	//unmarshal encrypted account
	var acc Store
	if err := json.Unmarshal([]byte(encryptedAccount), &acc); err != nil {
		return &KeyManager{}, err
	}

	//decrypt password with mnemonic
	pw, err := scrypt.DecryptCipherText(acc.Password, mnemonic)
	if err != nil {
		return &KeyManager{}, err
	}

	return OpenWithPassword(encryptedAccount, pw)

}

//Export the account
func (km KeyManager) Export(pw, pwConfirm string) (string, error) {

	//Exit if password's are not equal
	if pw != pwConfirm {
		return "", errors.New("password miss match")
	}

	//Marshal the keystore
	keyStore, err := km.keyStore.Marshal()
	if err != nil {
		return "", err
	}

	//encrypt key store with password
	encryptedKeyStore, err := scrypt.NewCipherText(string(keyStore), pw)
	if err != nil {
		return "", err
	}

	//encrypt password with mnemonic
	encryptedPassword, err := scrypt.NewCipherText(pw, km.keyStore.GetMnemonic().String())

	//Marshal account
	acc, err := json.Marshal(Store{
		Password:          encryptedPassword,
		EncryptedKeyStore: encryptedKeyStore,
		Version:           1,
	})

	return string(acc), err

}

//Get ethereum private key
func (km KeyManager) GetEthereumPrivateKey() (string, error) {
	return km.keyStore.GetKey(ethereumMigration.KeyStoreKey)
}

//Get ethereum address
func (km KeyManager) GetEthereumAddress() (string, error) {

	//Fetch ethereum private key
	privKey, err := km.GetEthereumPrivateKey()
	if err != nil {
		return "", err
	}

	//Parse hex private key
	priv, err := ethCrypto.HexToECDSA(privKey)
	if err != nil {
		return "", err
	}

	return ethCrypto.PubkeyToAddress(priv.PublicKey).String(), nil
}

func (km KeyManager) IdentityPrivateKey() (string, error) {
	return km.keyStore.GetKey(identity.Ed25519PrivateKey)
}

func (km KeyManager) IdentityPublicKey() (string, error) {
	return km.keyStore.GetKey(identity.Ed25519PublicKey)
}

func (km KeyManager) GetMnemonic() mnemonic.Mnemonic {
	return km.keyStore.GetMnemonic()
}

//Get the Mesh network private key (which is the identity ed25519 private key)
func (km KeyManager) MeshPrivateKey() (lp2pCrypto.PrivKey, error) {

	//Fetch private key
	priv, err := km.IdentityPrivateKey()
	if err != nil {
		return nil, err
	}
	privBytes, err := hex.DecodeString(priv)
	if err != nil {
		return nil, err
	}

	//Fetch public key
	pub, err := km.IdentityPublicKey()
	if err != nil {
		return nil, err
	}
	pubBytes, err := hex.DecodeString(pub)
	if err != nil {
		return nil, err
	}

	return lp2pCrypto.UnmarshalEd25519PrivateKey(append(privBytes, pubBytes...))

}

func (km KeyManager) GetEthereumPublicKey() (string, error) {

	//Fetch ethereum private key
	privKey, err := km.GetEthereumPrivateKey()
	if err != nil {
		return "", err
	}

	//Parse hex private key
	priv, err := ethCrypto.HexToECDSA(privKey)
	if err != nil {
		return "", err
	}

	// encode public key
	pubKey := ethCrypto.CompressPubkey(&priv.PublicKey)
	return hex.EncodeToString(pubKey), nil

}

//Sign data with identity key
func (km KeyManager) IdentitySign(data []byte) ([]byte, error) {

	//Fetch mesh private key
	pk, err := km.MeshPrivateKey()
	if err != nil {
		return nil, err
	}

	return pk.Sign(data)

}

//Sign data with ethereum private key
func (km KeyManager) EthereumSign(data [32]byte) ([]byte, error) {
	//Fetch ethereum private key
	privKey, err := km.GetEthereumPrivateKey()
	if err != nil {
		return nil, err
	}

	//Parse hex private key
	priv, err := ethCrypto.HexToECDSA(privKey)
	if err != nil {
		return nil, err
	}

	return ethCrypto.Sign(data[:], priv)
}

//Did the keystore change (happen after migration)
func (km KeyManager) WasMigrated() bool {
	return km.keyStore.WasMigrated()
}

func (km KeyManager) encryptionKey() (string, error) {
	return km.keyStore.GetKey(encryptionKeyMigration.BIP39Password)
}

func (km KeyManager) AESDecrypt(cipherText string) (string, error) {
	encryptionKey, err := km.encryptionKey()

	if err != nil {
		return "", err
	}

	return aes.Decrypt(cipherText, encryptionKey)
}

func (km KeyManager) AESEncrypt(plainText string) (string, error) {
	encryptionKey, err := km.encryptionKey()

	if err != nil {
		return "", err
	}

	return aes.Encrypt(plainText, encryptionKey)
}

func (km KeyManager) ChatIdKeyPair() (x3dh.KeyPair, error) {

	strPriv, err := km.keyStore.GetKey(chatMigration.MigrationPrivPrefix)
	if err != nil {
		return x3dh.KeyPair{}, err
	}
	rawPriv, err := hex.DecodeString(strPriv)
	if err != nil {
		return x3dh.KeyPair{}, err
	}

	strPub, err := km.keyStore.GetKey(chatMigration.MigrationPubPrefix)
	if err != nil {
		return x3dh.KeyPair{}, err
	}
	rawPub, err := hex.DecodeString(strPub)
	if err != nil {
		return x3dh.KeyPair{}, err
	}

	var (
		priv [32]byte
		pub  [32]byte
	)

	copy(priv[:], rawPriv[:32])
	copy(pub[:], rawPub[:32])

	return x3dh.KeyPair{
		PrivateKey: priv,
		PublicKey:  pub,
	}, nil
}

//Create new key manager from key store
func CreateFromKeyStore(ks ks.Store) *KeyManager {
	return &KeyManager{
		keyStore: ks,
	}
}
