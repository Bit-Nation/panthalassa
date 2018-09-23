package ethAddress

import (
	"testing"

	keyManager "github.com/Bit-Nation/panthalassa/keyManager"
	keyStore "github.com/Bit-Nation/panthalassa/keyStore"
	mnemonic "github.com/Bit-Nation/panthalassa/mnemonic"
	require "github.com/stretchr/testify/require"
	duktape "gopkg.in/olebedev/go-duktape.v3"
)

func TestModule_Register(t *testing.T) {

	// mnemonic
	mne, err := mnemonic.FromString("differ destroy head candy imitate barely wine ranch roof barrel sheriff blame umbrella visit sell green dress embark ramp cement rotate crawl session broom")
	require.Nil(t, err)

	// key store
	store, err := keyStore.NewFromMnemonic(mne)
	require.Nil(t, err)

	// open key manager form key store
	km := keyManager.CreateFromKeyStore(store)

	// create address module
	mod := New(km)

	vm := duktape.New()

	mod.Register(vm)

	v := vm.SafeToString(0)
	require.Equal(t, "0x748A6536dE0a8b1902f808233DD75ec4451cdFC6", v)

}
