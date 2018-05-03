package keyStore

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetKey(t *testing.T) {
	ks := KeyStore{
		keys: map[string]string{
			"key": "value",
		},
	}

	value, err := ks.GetKey("key")
	require.Equal(t, "value", value)
	require.Nil(t, err)

	value, err = ks.GetKey("i_do_not_exist_in_the_map")
	require.Equal(t, "key does not exist", err.Error())
}

func TestEthPrivateKeyValidationFunction(t *testing.T) {

	//Success test
	ks := KeyStore{
		mnemonic: "differ destroy head candy imitate barely wine ranch roof barrel sheriff blame umbrella visit sell green dress embark ramp cement rotate crawl session broom",
		keys: map[string]string{
			"eth_private_key": "eba47c97d7a6688d03e41b145d26090216c4468231bb46677553141f75222d5c",
		},
	}

	err := ethPrivateKeyValidation(ks)
	require.Equal(t, nil, err, "The error should be nil since we passed in the correct mnemonic")

	//Fail test
	ks = KeyStore{
		//We changed the first word which will result in an mismatch
		mnemonic: "destroy destroy head candy imitate barely wine ranch roof barrel sheriff blame umbrella visit sell green dress embark ramp cement rotate crawl session broom",
		keys: map[string]string{
			"eth_private_key": "eba47c97d7a6688d03e41b145d26090216c4468231bb46677553141f75222d5c",
		},
	}
	err = ethPrivateKeyValidation(ks)
	require.NotNil(t, err)
	require.Equal(t, "derivation mismatch - ethereum private key from storage and derived one doesn't match", err.Error())

}

//Test validation method
func TestValidateMethodOfKeyStore(t *testing.T) {

	ks := KeyStore{
		mnemonic: "differ destroy head candy imitate barely wine ranch roof barrel sheriff blame umbrella visit sell green dress embark ramp cement rotate crawl session broom",
		keys: map[string]string{
			"eth_private_key": "eba47c97d7a6688d03e41b145d26090216c4468231bb46677553141f75222d5c",
		},
		version: 1,
	}

	oldValidationRules := validationRules

	//Test validation of key store in case no rule set is present
	validationRules = map[uint8][]func(ks KeyStore) error{}
	err := ks.validate()
	require.Equal(t, "couldn't find validation rules", err.Error())

	//Test validation of key store if validation rules are present
	validationRules = oldValidationRules
	err = ks.validate()
	require.Nil(t, err)

	//Test validation of key store if validation set is present but not satisfied
	ks.mnemonic = "destroy destroy head candy imitate barely wine ranch roof barrel sheriff blame umbrella visit sell green dress embark ramp cement rotate crawl session broom" //Changed last word to house
	err = ks.validate()
	require.Equal(t, "derivation mismatch - ethereum private key from storage and derived one doesn't match", err.Error())
}

func TestJsonMarshalling(t *testing.T) {

	ks := KeyStore{
		mnemonic: "differ destroy head candy imitate barely wine ranch roof barrel sheriff blame umbrella visit sell green dress embark ramp cement rotate crawl session broom",
		keys: map[string]string{
			"eth_private_key": "eba47c97d7a6688d03e41b145d26090216c4468231bb46677553141f75222d5c",
		},
		version: 1,
	}
	oldValidationRules := validationRules

	//Test that marshal fails when validation fails
	validationRules = map[uint8][]func(ks KeyStore) error{
		1: []func(store KeyStore) error{
			func(store KeyStore) error {
				return errors.New("validation failed")
			},
		},
	}
	_, err := ks.Marshal()
	require.Equal(t, "validation failed", err.Error())
	validationRules = oldValidationRules

	//Test successful json marshal
	expectedKeyStoreExport := []byte(`{"mnemonic":"differ destroy head candy imitate barely wine ranch roof barrel sheriff blame umbrella visit sell green dress embark ramp cement rotate crawl session broom","keys":{"eth_private_key":"eba47c97d7a6688d03e41b145d26090216c4468231bb46677553141f75222d5c"},"version":1}`)
	json, err := ks.Marshal()
	require.Nil(t, err)
	require.Equal(t, expectedKeyStoreExport, json)

}

func TestKeyStoreJsonImport(t *testing.T) {

	//import a valid json key store
	jsonKeyStore := `{"mnemonic":"differ destroy head candy imitate barely wine ranch roof barrel sheriff blame umbrella visit sell green dress embark ramp cement rotate crawl session broom","keys":{"eth_private_key":"eba47c97d7a6688d03e41b145d26090216c4468231bb46677553141f75222d5c"},"version":1}`

	ks, err := FromJson(jsonKeyStore)
	require.Nil(t, err)

	require.Equal(t, "differ destroy head candy imitate barely wine ranch roof barrel sheriff blame umbrella visit sell green dress embark ramp cement rotate crawl session broom", ks.mnemonic)
	require.Equal(t, ks.keys["eth_private_key"], "eba47c97d7a6688d03e41b145d26090216c4468231bb46677553141f75222d5c")
	require.Equal(t, ks.version, uint8(1))

	//import invalid json keystore
	jsonKeyStore = `{}`
	_, err = FromJson(jsonKeyStore)
	require.Equal(t, "couldn't find validation rules", err.Error())

}

func TestNewKeyStoreFactory(t *testing.T) {

	newMnemonic = func() (string, error) {
		return "differ destroy head candy imitate barely wine ranch roof barrel sheriff blame umbrella visit sell green dress embark ramp cement rotate crawl session broom", nil
	}

	ks, err := NewKeyStoreFactory()
	require.Nil(t, err)

	require.Equal(t, "differ destroy head candy imitate barely wine ranch roof barrel sheriff blame umbrella visit sell green dress embark ramp cement rotate crawl session broom", ks.mnemonic)
	require.Equal(t, ks.keys["eth_private_key"], "eba47c97d7a6688d03e41b145d26090216c4468231bb46677553141f75222d5c")
	require.Equal(t, ks.version, uint8(1))

}

func TestGetMnemonic(t *testing.T) {

	ks := KeyStore{
		mnemonic: "differ destroy head candy imitate barely wine ranch roof barrel sheriff blame umbrella visit sell green dress embark ramp cement rotate crawl session broom",
	}

	require.Equal(t, "differ destroy head candy imitate barely wine ranch roof barrel sheriff blame umbrella visit sell green dress embark ramp cement rotate crawl session broom", ks.GetMnemonic())

}
