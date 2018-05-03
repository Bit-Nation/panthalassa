package keys

import "github.com/Bit-Nation/panthalassa/mnemonic"

type Migration interface {
	Up(mnemonic mnemonic.Mnemonic, keys map[string]string) (map[string]string, error)
}