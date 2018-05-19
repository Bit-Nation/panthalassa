package panthalassa

import (
	api "github.com/Bit-Nation/panthalassa/api/device"
	deviceApi "github.com/Bit-Nation/panthalassa/api/device"
	keyManager "github.com/Bit-Nation/panthalassa/keyManager"
	mesh "github.com/Bit-Nation/panthalassa/mesh"
)

type Panthalassa struct {
	km        *keyManager.KeyManager
	upStream  api.UpStream
	deviceApi *deviceApi.Api
	mesh      *mesh.Network
}

//Stop the panthalassa instance
//this becomes interesting when we start
//to use the mesh network
func (p *Panthalassa) Stop() error {
	return p.mesh.Close()
}

//Export account with the given password
func (p *Panthalassa) Export(pw, pwConfirm string) (string, error) {
	return p.km.Export(pw, pwConfirm)
}
