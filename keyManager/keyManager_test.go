package keyManager

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/scrypt"
	"testing"

	aes "github.com/Bit-Nation/panthalassa/crypto/aes"
	panthScrypt "github.com/Bit-Nation/panthalassa/crypto/scrypt"
	keyStore "github.com/Bit-Nation/panthalassa/keyStore"
	identity "github.com/Bit-Nation/panthalassa/keyStore/migration/identity"
	mnemonic "github.com/Bit-Nation/panthalassa/mnemonic"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	require "github.com/stretchr/testify/require"
)

//Test the Create from function
func TestCreateFromKeyStore(t *testing.T) {

	//mnemonic
	mn, err := mnemonic.New()
	require.Nil(t, err)

	//create keyStore
	ks, err := keyStore.NewFromMnemonic(mn)
	require.Nil(t, err)

	km := CreateFromKeyStore(ks)

	require.Equal(t, km.keyStore, ks)
}

func TestExportFunction(t *testing.T) {

	//create key storage
	jsonKeyStore := `{"mnemonic":"differ destroy head candy imitate barely wine ranch roof barrel sheriff blame umbrella visit sell green dress embark ramp cement rotate crawl session broom","keys":{"chat_identity_curve25519_private_key":"70bcdb281ab3cc1dc75199c33a0edec43fcfe1d70ee2fd11e4821c38a688186c","chat_identity_curve25519_public_key":"1b276c51c849b244a7c40814769c9ea71caad17516aabc1270c8bd2bc096ef45","ed_25519_private_key":"9d426d0eb4170529672df197454bc77cc36cb341c872bcee0bece79ac893b34a8c5de2e7d099b881ed6214f8add6cbba2a84f57546b7f0a6d39197c904529f3f","ed_25519_public_key":"8c5de2e7d099b881ed6214f8add6cbba2a84f57546b7f0a6d39197c904529f3f","encryption_key":"7dc02d78d98fff23d1f4500e4c8742fb26ad233db2d421d5bcb44306a2bb69e2","ethereum_private_key":"eba47c97d7a6688d03e41b145d26090216c4468231bb46677553141f75222d5c"},"version":1}`
	ks, err := keyStore.UnmarshalStore(jsonKeyStore)
	require.Nil(t, err)

	//create key manager
	km := CreateFromKeyStore(ks)

	//Export the key storage via the key manager
	//The export should be encrypted
	cipherText, err := km.Export("my_password", "my_password")
	require.Nil(t, err)

	//Decrypt the exported encrypted key storage
	km, err = OpenWithPassword(cipherText, "my_password")
	require.Nil(t, err)

	jsonKs, err := km.keyStore.Marshal()
	require.Nil(t, err)

	require.Equal(t, jsonKeyStore, string(jsonKs))
}

func TestOpenWithMnemonic(t *testing.T) {

	//create key storage
	jsonKeyStore := `{"mnemonic":"differ destroy head candy imitate barely wine ranch roof barrel sheriff blame umbrella visit sell green dress embark ramp cement rotate crawl session broom","keys":{"chat_identity_curve25519_private_key":"70bcdb281ab3cc1dc75199c33a0edec43fcfe1d70ee2fd11e4821c38a688186c","chat_identity_curve25519_public_key":"1b276c51c849b244a7c40814769c9ea71caad17516aabc1270c8bd2bc096ef45","ed_25519_private_key":"9d426d0eb4170529672df197454bc77cc36cb341c872bcee0bece79ac893b34a8c5de2e7d099b881ed6214f8add6cbba2a84f57546b7f0a6d39197c904529f3f","ed_25519_public_key":"8c5de2e7d099b881ed6214f8add6cbba2a84f57546b7f0a6d39197c904529f3f","encryption_key":"7dc02d78d98fff23d1f4500e4c8742fb26ad233db2d421d5bcb44306a2bb69e2","ethereum_private_key":"eba47c97d7a6688d03e41b145d26090216c4468231bb46677553141f75222d5c"},"version":1}`
	ks, err := keyStore.UnmarshalStore(jsonKeyStore)
	require.Nil(t, err)

	//create key manager
	km := CreateFromKeyStore(ks)

	//Export the key storage via the key manager
	//The export should be encrypted
	cipherText, err := km.Export("my_password", "my_password")
	require.Nil(t, err)

	//Decrypt the exported encrypted key storage
	km, err = OpenWithMnemonic(cipherText, "differ destroy head candy imitate barely wine ranch roof barrel sheriff blame umbrella visit sell green dress embark ramp cement rotate crawl session broom")
	require.Nil(t, err)

	jsonKs, err := km.keyStore.Marshal()
	require.Nil(t, err)

	require.Equal(t, jsonKeyStore, string(jsonKs))

}

func TestGetMnemonic(t *testing.T) {

	//create key storage
	jsonKeyStore := `{"mnemonic":"differ destroy head candy imitate barely wine ranch roof barrel sheriff blame umbrella visit sell green dress embark ramp cement rotate crawl session broom","keys":{"ethereum_private_key":"eba47c97d7a6688d03e41b145d26090216c4468231bb46677553141f75222d5c"},"version":1}`
	ks, err := keyStore.UnmarshalStore(jsonKeyStore)
	require.Nil(t, err)

	//key manager
	km := CreateFromKeyStore(ks)

	//Get address
	mne := km.GetMnemonic()
	require.Equal(t, "differ destroy head candy imitate barely wine ranch roof barrel sheriff blame umbrella visit sell green dress embark ramp cement rotate crawl session broom", mne.String())

}

func TestGetAddressFromPrivateKey(t *testing.T) {

	//create key storage
	jsonKeyStore := `{"mnemonic":"differ destroy head candy imitate barely wine ranch roof barrel sheriff blame umbrella visit sell green dress embark ramp cement rotate crawl session broom","keys":{"ethereum_private_key":"eba47c97d7a6688d03e41b145d26090216c4468231bb46677553141f75222d5c"},"version":1}`
	ks, err := keyStore.UnmarshalStore(jsonKeyStore)
	require.Nil(t, err)

	//key manager
	km := CreateFromKeyStore(ks)

	//Get address
	addr, err := km.GetEthereumAddress()
	require.Nil(t, err)
	require.Equal(t, "0x748A6536dE0a8b1902f808233DD75ec4451cdFC6", addr)

}

func TestGetLibP2PPrivateKey(t *testing.T) {

	//create key storage
	jsonKeyStore := `{"mnemonic":"differ destroy head candy imitate barely wine ranch roof barrel sheriff blame umbrella visit sell green dress embark ramp cement rotate crawl session broom","keys":{"ethereum_private_key":"eba47c97d7a6688d03e41b145d26090216c4468231bb46677553141f75222d5c"},"version":1}`
	ks, err := keyStore.UnmarshalStore(jsonKeyStore)
	require.Nil(t, err)

	//key manager
	km := CreateFromKeyStore(ks)

	meshPriv, err := km.MeshPrivateKey()
	require.Nil(t, err)

	edPriv, err := ks.GetKey(identity.Ed25519PrivateKey)
	require.Nil(t, err)
	edPrivBytes, err := hex.DecodeString(edPriv)
	require.Nil(t, err)

	edPub, err := ks.GetKey(identity.Ed25519PublicKey)
	require.Nil(t, err)
	edPubBytes, err := hex.DecodeString(edPub)
	require.Nil(t, err)

	//Check if private key does match
	meshPrivBytes, err := meshPriv.Bytes()
	require.Nil(t, err)
	combinedKey := append(edPrivBytes, edPubBytes...)
	preFix, err := hex.DecodeString("08011260")
	require.Nil(t, err)
	require.Equal(t, hex.EncodeToString(append(preFix, combinedKey...)), hex.EncodeToString(meshPrivBytes))
}

func TestGetEthereumPublicKey(t *testing.T) {

	//create key storage
	jsonKeyStore := `{"mnemonic":"differ destroy head candy imitate barely wine ranch roof barrel sheriff blame umbrella visit sell green dress embark ramp cement rotate crawl session broom","keys":{"ethereum_private_key":"eba47c97d7a6688d03e41b145d26090216c4468231bb46677553141f75222d5c"},"version":1}`
	ks, err := keyStore.UnmarshalStore(jsonKeyStore)
	require.Nil(t, err)

	km := CreateFromKeyStore(ks)

	ethPublicKey, err := km.GetEthereumPublicKey()
	require.Nil(t, err)

	require.Equal(t, "032b6a023528114bdf34718260a18b520def9705ea2b3c0ec41e160204a5fa8493", ethPublicKey)

}

func TestGetEthereumAddress(t *testing.T) {

	//create key storage
	jsonKeyStore := `{"mnemonic":"differ destroy head candy imitate barely wine ranch roof barrel sheriff blame umbrella visit sell green dress embark ramp cement rotate crawl session broom","keys":{"ethereum_private_key":"eba47c97d7a6688d03e41b145d26090216c4468231bb46677553141f75222d5c"},"version":1}`
	ks, err := keyStore.UnmarshalStore(jsonKeyStore)
	require.Nil(t, err)

	km := CreateFromKeyStore(ks)

	address, err := km.GetEthereumAddress()
	require.Nil(t, err)

	require.Equal(t, "0x748A6536dE0a8b1902f808233DD75ec4451cdFC6", address)

}

func TestKeyManager_EthereumSign(t *testing.T) {

	//create key storage
	jsonKeyStore := `{"mnemonic":"differ destroy head candy imitate barely wine ranch roof barrel sheriff blame umbrella visit sell green dress embark ramp cement rotate crawl session broom","keys":{"ed_25519_private_key":"9d426d0eb4170529672df197454bc77cc36cb341c872bcee0bece79ac893b34a8c5de2e7d099b881ed6214f8add6cbba2a84f57546b7f0a6d39197c904529f3f","ed_25519_public_key":"8c5de2e7d099b881ed6214f8add6cbba2a84f57546b7f0a6d39197c904529f3f","ethereum_private_key":"eba47c97d7a6688d03e41b145d26090216c4468231bb46677553141f75222d5c"},"version":1}`
	ks, err := keyStore.UnmarshalStore(jsonKeyStore)
	require.Nil(t, err)

	// create key manager
	km := CreateFromKeyStore(ks)

	// sample hash
	hash := sha256.Sum256([]byte("hi"))

	// sign data
	signature, err := km.EthereumSign(hash)
	require.Nil(t, err)

	pubKeyStr, err := km.GetEthereumPublicKey()
	require.Nil(t, err)
	rawPubKey, err := hex.DecodeString(pubKeyStr)
	require.Nil(t, err)

	// should be true since the signature should be valid
	require.True(t, ethCrypto.VerifySignature(rawPubKey, hash[:], signature[:64]))

}

func TestKeyManager_IdentitySign(t *testing.T) {

	//create key storage
	jsonKeyStore := `{"mnemonic":"differ destroy head candy imitate barely wine ranch roof barrel sheriff blame umbrella visit sell green dress embark ramp cement rotate crawl session broom","keys":{"ed_25519_private_key":"9d426d0eb4170529672df197454bc77cc36cb341c872bcee0bece79ac893b34a8c5de2e7d099b881ed6214f8add6cbba2a84f57546b7f0a6d39197c904529f3f","ed_25519_public_key":"8c5de2e7d099b881ed6214f8add6cbba2a84f57546b7f0a6d39197c904529f3f","ethereum_private_key":"eba47c97d7a6688d03e41b145d26090216c4468231bb46677553141f75222d5c"},"version":1}`
	ks, err := keyStore.UnmarshalStore(jsonKeyStore)
	require.Nil(t, err)

	// create key manager
	km := CreateFromKeyStore(ks)

	// sign data
	signedData, err := km.IdentitySign([]byte("hi"))
	require.Nil(t, err)

	// fetch public key
	idPubStr, err := km.IdentityPublicKey()
	require.Nil(t, err)
	pub, err := hex.DecodeString(idPubStr)
	require.Nil(t, err)

	// signature should be valid
	require.True(t, ed25519.Verify(pub, []byte("hi"), signedData))

}

func TestKeyManager_AESEncryptDecrypt(t *testing.T) {

	//create key storage
	jsonKeyStore := `{"mnemonic":"differ destroy head candy imitate barely wine ranch roof barrel sheriff blame umbrella visit sell green dress embark ramp cement rotate crawl session broom","keys":{"ed_25519_private_key":"9d426d0eb4170529672df197454bc77cc36cb341c872bcee0bece79ac893b34a8c5de2e7d099b881ed6214f8add6cbba2a84f57546b7f0a6d39197c904529f3f","ed_25519_public_key":"8c5de2e7d099b881ed6214f8add6cbba2a84f57546b7f0a6d39197c904529f3f","ethereum_private_key":"eba47c97d7a6688d03e41b145d26090216c4468231bb46677553141f75222d5c"},"version":1}`
	ks, err := keyStore.UnmarshalStore(jsonKeyStore)
	require.Nil(t, err)

	// create key manager
	km := CreateFromKeyStore(ks)

	// encrypt the cipher text
	cipherText, err := km.AESEncrypt([]byte("hi"))
	require.Nil(t, err)

	// decrypt the cipher text
	plain, err := km.AESDecrypt(cipherText)
	require.Nil(t, err)
	require.Equal(t, "hi", string(plain))
}

func TestV2KeyManager(t *testing.T) {

	//create key storage
	jsonKeyStore := `{"mnemonic":"differ destroy head candy imitate barely wine ranch roof barrel sheriff blame umbrella visit sell green dress embark ramp cement rotate crawl session broom","keys":{"ed_25519_private_key":"9d426d0eb4170529672df197454bc77cc36cb341c872bcee0bece79ac893b34a8c5de2e7d099b881ed6214f8add6cbba2a84f57546b7f0a6d39197c904529f3f","ed_25519_public_key":"8c5de2e7d099b881ed6214f8add6cbba2a84f57546b7f0a6d39197c904529f3f","ethereum_private_key":"eba47c97d7a6688d03e41b145d26090216c4468231bb46677553141f75222d5c"},"version":1}`
	ks, err := keyStore.UnmarshalStore(jsonKeyStore)
	require.Nil(t, err)

	// create key manager
	km := CreateFromKeyStore(ks)

	// test if export works
	s, err := km.Export("pw", "pw")
	require.Nil(t, err)

	// must be version 2 from now on
	require.Equal(t, uint8(2), s.Version)

}

func TestMigration(t *testing.T) {

	pw := "my_password"

	// scrypt params
	const n = 16384
	const r = 8
	const p = 1
	const keyLength = 32

	// create salt
	salt := make([]byte, 50)
	_, err := rand.Read(salt)
	require.Nil(t, err)

	// mnemonic
	mne, err := mnemonic.FromString("differ destroy head candy imitate barely wine ranch roof barrel sheriff blame umbrella visit sell green dress embark ramp cement rotate crawl session broom")
	require.Nil(t, err)

	// create key store from mnemonic
	s, err := keyStore.NewFromMnemonic(mne)
	require.Nil(t, err)

	// fetch private key
	ethPrivateKey, err := s.GetKey("ethereum_private_key")
	require.Nil(t, err)

	// export key store
	strKeyStore, err := s.Marshal()
	require.Nil(t, err)

	// create old version of key manager
	// version 1 used to use AES CFB (scrypt cipher v1)
	key, err := scrypt.Key([]byte(pw), salt, n, r, p, keyLength)
	require.Nil(t, err)
	var encryptionSecret = aes.Secret{}
	copy(encryptionSecret[:], key[:])

	ct, err := aes.CFBEncrypt(strKeyStore, encryptionSecret)
	require.Nil(t, err)

	// this is a valid panthalassa scrypt cipher text
	// which can be used by the key manager
	// (key manager pass it ot our scrypt package)
	encryptedKeyStore := panthScrypt.CipherText{
		CipherText: ct,
		ScryptKey: panthScrypt.Key{
			N:      n,
			P:      p,
			R:      r,
			KeyLen: keyLength,
			Salt:   salt,
		},
		Version: uint8(0),
	}

	// create key manager store
	str := Store{
		EncryptedKeyStore: encryptedKeyStore,
		Version:           1,
	}

	// open the store we just created
	// AES CFB should be used to decrypt
	// the key store
	km, err := OpenWithPassword(str, pw)
	require.Nil(t, err)

	// when we export it we should now
	// have a key manager key store
	// of version 2 since we use AES CTR
	// from now on.
	migratedStore, err := km.Export(pw, pw)
	require.Nil(t, err)

	require.Equal(t, uint8(2), migratedStore.Version)
	require.Equal(t, uint8(2), migratedStore.EncryptedKeyStore.CipherText.Version)

	// open the migrate store again
	// to make sure decryption works fine
	km, err = OpenWithPassword(migratedStore, pw)
	require.Nil(t, err)

	// get ethereum private key
	pk, err := km.GetEthereumPrivateKey()
	require.Nil(t, err)

	// make sure that the private key is still the same
	// after the migration from CFB to CTR
	require.Equal(t, ethPrivateKey, pk)

}
