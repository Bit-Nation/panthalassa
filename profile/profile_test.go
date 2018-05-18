package profile

import (
	"testing"

	km "github.com/Bit-Nation/panthalassa/keyManager"
	ks "github.com/Bit-Nation/panthalassa/keyStore"
	mnemonic "github.com/Bit-Nation/panthalassa/mnemonic"
	require "github.com/stretchr/testify/require"
)

func TestSignProfile(t *testing.T) {

	// create test mnemonic
	mne, err := mnemonic.New()
	require.Nil(t, err)

	// create key store
	store, err := ks.NewFromMnemonic(mne)
	require.Nil(t, err)

	// open key manger with created keystore
	keyManager := km.CreateFromKeyStore(store)

	// create profile
	prof, err := SignProfile("Florian", "Germany", "base64", *keyManager)
	require.Nil(t, err)

	// validate profile
	valid, err := prof.SignaturesValid()
	require.Nil(t, err)
	require.True(t, valid)

}

/*
func TestCalculateProfile(t *testing.T) {

	// create test mnemonic
	mne, err := mnemonic.New()
	require.Nil(t, err)

	// create key store
	store, err := ks.NewFromMnemonic(mne)
	require.Nil(t, err)

	// open key manger with created keystore
	keyManager := km.CreateFromKeyStore(store)

	// create panthalassa instance
	p := panthalassa{
		km: keyManager,
	}

	// calculate the profile. The profile will contain metadata and public key's
	// signatures of the public key's are attached to verify the integrity of data
	res, err := p.CalculateProfile("Florian", "Earth", "base64..")
	require.Nil(t, err)

	// unmarshal account to check if it work's like we expect it to work
	var prof profile
	require.Nil(t, json.Unmarshal([]byte(res), &prof))

	// basic check's
	require.Equal(t, "Florian", prof.Information.Name)
	require.Equal(t, "Earth", prof.Information.Location)
	require.Equal(t, "base64..", prof.Information.Image)

	// id pub key as string
	idPubKeyStr, err := p.km.IdentityPublicKey()
	require.Nil(t, err)

	// unmarshal id pub key
	idPubKey, err := hex.DecodeString(idPubKeyStr)
	require.Nil(t, err)

	// eth pub key
	ethPubKeyStr, err := p.km.GetEthereumPublicKey()
	require.Nil(t, err)

	// unmarshal id pub key
	ethPubKey, err := hex.DecodeString(ethPubKeyStr)
	require.Nil(t, err)

	// concat profile information
	dataToSign := []byte("Florian")
	dataToSign = append(dataToSign, []byte("Earth")...)
	dataToSign = append(dataToSign, []byte("base64..")...)
	dataToSign = append(dataToSign, idPubKey...)
	dataToSign = append(dataToSign, ethPubKey...)
	version := make([]byte, 2)
	binary.LittleEndian.PutUint16(version, profileVersion)
	dataToSign = append(dataToSign, version...)

	// hash of profile data
	dataHash := sha256.Sum256(dataToSign)

	// d
	edIdKey, err := lp2pCrypto.UnmarshalEd25519PublicKey(idPubKey)
	idKeySignature, err := hex.DecodeString(prof.Signatures.IdentityKey)
	fmt.Println(edIdKey.Verify(dataHash[:], idKeySignature))

}
*/
