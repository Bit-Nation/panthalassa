package panthalassa

import (
	"errors"

	deviceApi "github.com/Bit-Nation/panthalassa/api/device"
	keyManager "github.com/Bit-Nation/panthalassa/keyManager"
	mesh "github.com/Bit-Nation/panthalassa/mesh"
	log "github.com/ipfs/go-log"
)

var panthalassaInstance *panthalassa
var logger = log.Logger("panthalassa")

type UpStream interface {
	Send(data string)
}

//Create a new panthalassa instance
func Start(accountStore, password, rendezvousKey string, upStream UpStream) error {

	//Exit if instance was already created and not stopped
	if panthalassaInstance != nil {
		return errors.New("call stop first in order to create a new panthalassa instance")
	}

	//Create key manager
	km, err := keyManager.OpenWithPassword(accountStore, password)
	if err != nil {
		return err
	}

	//Mesh network
	pk, err := km.MeshPrivateKey()
	if err != nil {
		return err
	}

	m, errReporter, err := mesh.New(pk, rendezvousKey)
	if err != nil {
		return err
	}
	//Report error's from mesh network to current logger
	go func() {
		for {
			select {
			case err := <-errReporter:
				logger.Error(err)
			}
		}
	}()

	//Create panthalassa instance
	panthalassaInstance = &panthalassa{
		km:        km,
		upStream:  upStream,
		deviceApi: deviceApi.New(upStream),
		mesh:      m,
	}

	return nil

}

//Create a new panthalassa instance with the mnemonic
func StartFromMnemonic(accountStore, mnemonic string) error {

	if panthalassaInstance != nil {
		return errors.New("call stop first in order to create a new panthalassa instance")
	}

	//Create key manager
	km, err := keyManager.OpenWithMnemonic(accountStore, mnemonic)
	if err != nil {
		return err
	}

	//Create panthalassa instance
	panthalassaInstance = &panthalassa{
		km: km,
	}

	return nil

}

//Eth Private key
func EthPrivateKey() (string, error) {

	if panthalassaInstance == nil {
		return "", errors.New("you have to start panthalassa")
	}

	return panthalassaInstance.km.GetEthereumPrivateKey()

}

func EthAddress() (string, error) {
	if panthalassaInstance == nil {
		return "", errors.New("you have to start panthalassa")
	}

	return panthalassaInstance.km.GetEthereumAddress()
}

func SendResponse(id uint32, data string) error {

	if panthalassaInstance == nil {
		return errors.New("you have to start panthalassa")
	}

	return panthalassaInstance.deviceApi.Receive(id, data)
}

//Export the current account store with given password
func ExportAccountStore(pw, pwConfirm string) (string, error) {

	if panthalassaInstance == nil {
		return "", errors.New("you have to start panthalassa")
	}

	return panthalassaInstance.Export(pw, pwConfirm)

}

func IdentityPublicKey() (string, error) {

	if panthalassaInstance == nil {
		return "", errors.New("you have to start panthalassa")
	}

	return panthalassaInstance.km.IdentityPublicKey()
}

//Stop panthalassa
func Stop() error {

	//Exit if not started
	if panthalassaInstance == nil {
		return errors.New("you have to start panthalassa in order to stop it")
	}

	//Stop panthalassa
	err := panthalassaInstance.Stop()
	if err != nil {
		return err
	}

	//Reset singleton
	panthalassaInstance = nil

	return nil
}
