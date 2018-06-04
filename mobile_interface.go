package panthalassa

import (
	"encoding/json"
	"errors"

	deviceApi "github.com/Bit-Nation/panthalassa/api/device"
	"github.com/Bit-Nation/panthalassa/chat"
	keyManager "github.com/Bit-Nation/panthalassa/keyManager"
	mesh "github.com/Bit-Nation/panthalassa/mesh"
	profile "github.com/Bit-Nation/panthalassa/profile"
	log "github.com/ipfs/go-log"
	valid "gopkg.in/asaskevich/govalidator.v4"
)

var panthalassaInstance *Panthalassa
var logger = log.Logger("panthalassa")

type UpStream interface {
	Send(data string)
}

type StartConfig struct {
	EncryptedKeyManager string `valid:"required"`
	RendezvousKey       string `valid:"required"`
	SignedProfile       string `valid:"required"`
	ChatHTTPEndpoint    string `valid:"required"`
	ChatWSEndpoint      string `valid:"required"`
	ChatAccessToken     string `valid:"required"`
}

// create a new panthalassa instance
func start(km *keyManager.KeyManager, chatKeyStore PangeaKeyStoreDBInterface, config StartConfig, client UpStream) error {

	//Exit if instance was already created and not stopped
	if panthalassaInstance != nil {
		return errors.New("call stop first in order to create a new panthalassa instance")
	}

	//Mesh network
	pk, err := km.MeshPrivateKey()
	if err != nil {
		return err
	}

	// device api
	api := deviceApi.New(client)

	m, errReporter, err := mesh.New(pk, api, config.RendezvousKey, config.SignedProfile)
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

	chatKeyPair, err := km.ChatIdKeyPair()
	if err != nil {
		return err
	}

	// @Todo implement the key store and the client
	c, err := chat.New(chatKeyPair, km, nil, nil)
	if err != nil {
		return err
	}

	//Create panthalassa instance
	panthalassaInstance = &Panthalassa{
		km:        km,
		upStream:  client,
		deviceApi: api,
		mesh:      m,
		chat:      &c,
	}

	return nil

}

// start panthalassa
func Start(config string, password string, chatKeyStore PangeaKeyStoreDBInterface, client UpStream) error {

	// unmarshal config
	var c StartConfig
	if err := json.Unmarshal([]byte(config), &c); err != nil {
		return err
	}

	// validate config
	_, err := valid.ValidateStruct(c)
	if err != nil {
		return err
	}

	// open key manager with password
	km, err := keyManager.OpenWithPassword(c.EncryptedKeyManager, password)
	if err != nil {
		return err
	}

	return start(km, chatKeyStore, c, client)
}

// create a new panthalassa instance with the mnemonic
func StartFromMnemonic(config string, mnemonic string, chatKeyStore PangeaKeyStoreDBInterface, client UpStream) error {

	// unmarshal config
	var c StartConfig
	if err := json.Unmarshal([]byte(config), &c); err != nil {
		return err
	}

	// validate config
	_, err := valid.ValidateStruct(c)
	if err != nil {
		return err
	}

	// create key manager
	km, err := keyManager.OpenWithMnemonic(c.EncryptedKeyManager, mnemonic)
	if err != nil {
		return err
	}

	// create panthalassa instance
	return start(km, chatKeyStore, c, client)

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

func SendResponse(id string, data string) error {

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

func GetMnemonic() (string, error) {

	if panthalassaInstance == nil {
		return "", errors.New("you have to start panthalassa")
	}

	return panthalassaInstance.km.GetMnemonic().String(), nil
}

func SignProfile(name, location, image string) (string, error) {

	if panthalassaInstance == nil {
		return "", errors.New("you have to start panthalassa")
	}

	p, err := profile.SignProfile(name, location, image, *panthalassaInstance.km)
	if err != nil {
		return "", err
	}

	rawProfile, err := p.Marshal()
	if err != nil {
		return "", err
	}

	return string(rawProfile), nil

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

// fetch the identity public key of the
func GetIdentityPublicKey() (string, error) {

	//Exit if not started
	if panthalassaInstance == nil {
		return "", errors.New("you have to start panthalassa first")
	}

	return panthalassaInstance.km.IdentityPublicKey()

}
