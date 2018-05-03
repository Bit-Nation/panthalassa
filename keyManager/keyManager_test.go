package keyManager

import (
	keyStore "github.com/Bit-Nation/panthalassa/keyStore"
	"github.com/Bit-Nation/panthalassa/mnemonic"
	"github.com/stretchr/testify/require"
	"testing"
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
	jsonKeyStore := `{"mnemonic":"differ destroy head candy imitate barely wine ranch roof barrel sheriff blame umbrella visit sell green dress embark ramp cement rotate crawl session broom","keys":{"ethereum_private_key":"eba47c97d7a6688d03e41b145d26090216c4468231bb46677553141f75222d5c"},"version":1}`
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
	jsonKeyStore := `{"mnemonic":"differ destroy head candy imitate barely wine ranch roof barrel sheriff blame umbrella visit sell green dress embark ramp cement rotate crawl session broom","keys":{"ethereum_private_key":"eba47c97d7a6688d03e41b145d26090216c4468231bb46677553141f75222d5c"},"version":1}`
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
