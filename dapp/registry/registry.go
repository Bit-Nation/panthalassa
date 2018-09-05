package registry

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	api "github.com/Bit-Nation/panthalassa/api"
	dapp "github.com/Bit-Nation/panthalassa/dapp"
	module "github.com/Bit-Nation/panthalassa/dapp/module"
	ethAddrMod "github.com/Bit-Nation/panthalassa/dapp/module/ethAddress"
	loggerMod "github.com/Bit-Nation/panthalassa/dapp/module/logger"
	messageModule "github.com/Bit-Nation/panthalassa/dapp/module/message"
	modalMod "github.com/Bit-Nation/panthalassa/dapp/module/modal"
	randBytes "github.com/Bit-Nation/panthalassa/dapp/module/randBytes"
	renderDApp "github.com/Bit-Nation/panthalassa/dapp/module/renderer/dapp"
	renderMsg "github.com/Bit-Nation/panthalassa/dapp/module/renderer/message"
	sendEthTxMod "github.com/Bit-Nation/panthalassa/dapp/module/sendEthTx"
	uuidv4Mod "github.com/Bit-Nation/panthalassa/dapp/module/uuidv4"
	db "github.com/Bit-Nation/panthalassa/db"
	keyManager "github.com/Bit-Nation/panthalassa/keyManager"
	storm "github.com/asdine/storm"
	log "github.com/ipfs/go-log"
	host "github.com/libp2p/go-libp2p-host"
	net "github.com/libp2p/go-libp2p-net"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	ma "github.com/multiformats/go-multiaddr"
	golog "github.com/op/go-logging"
	ed25519 "golang.org/x/crypto/ed25519"
)

var logger = log.Logger("dapp - registry")

type fetchDAppChanStr struct {
	signingKey ed25519.PublicKey
	respChan   chan *dapp.DApp
}

type addDevStreamChanStr struct {
	signingKey ed25519.PublicKey
	stream     net.Stream
}

type fetchDAppStreamStr struct {
	signingKey ed25519.PublicKey
	respChan   chan net.Stream
}

// keep track of all running DApps
type Registry struct {
	host               host.Host
	closeChan          chan *dapp.Data
	conf               Config
	api                *api.API
	km                 *keyManager.KeyManager
	dAppDB             dapp.Storage
	msgDB              db.ChatStorage
	db                 *storm.DB
	addDAppChan        chan *dapp.DApp
	fetchDAppChan      chan fetchDAppChanStr
	addDevStreamChan   chan addDevStreamChanStr
	fetchDevStreamChan chan fetchDAppStreamStr
}

type Config struct {
	EthWSEndpoint string
}

// create new dApp registry
func NewDAppRegistry(h host.Host, conf Config, api *api.API, km *keyManager.KeyManager, dAppDB dapp.Storage, msgDB db.ChatStorage) (*Registry, error) {

	r := &Registry{
		host:               h,
		closeChan:          make(chan *dapp.Data),
		conf:               conf,
		api:                api,
		km:                 km,
		dAppDB:             dAppDB,
		msgDB:              msgDB,
		addDAppChan:        make(chan *dapp.DApp),
		fetchDAppChan:      make(chan fetchDAppChanStr),
		addDevStreamChan:   make(chan addDevStreamChanStr),
		fetchDevStreamChan: make(chan fetchDAppStreamStr),
	}

	// load all default DApps
	dApps := []dapp.RawData{}
	if err := json.Unmarshal([]byte(rawDApps), &dApps); err != nil {
		return nil, err
	}
	for _, dAppData := range dApps {
		dApp, err := dapp.ParseJsonToData(dAppData)
		if err != nil {
			return nil, err
		}
		if err := dAppDB.SaveDApp(dApp); err != nil {
			return nil, err
		}
	}

	// add worker to remove DApps
	go func() {

		dAppInstances := map[string]*dapp.DApp{}
		streams := map[string]net.Stream{}

		for {
			select {
			// remove DApp from state
			case cc := <-r.closeChan:
				delete(dAppInstances, hex.EncodeToString(cc.UsedSigningKey))
				// @todo send signal to client that this app was shut down
			// fetch dApp from state
			case dAppFetch := <-r.fetchDAppChan:
				dApp, exist := dAppInstances[hex.EncodeToString(dAppFetch.signingKey)]
				if !exist {
					dAppFetch.respChan <- nil
					continue
				}
				dAppFetch.respChan <- dApp
			// add dev stream to state
			case addStream := <-r.addDevStreamChan:
				streams[hex.EncodeToString(addStream.signingKey)] = addStream.stream
			// fetch dev stream from state
			case fetchDevStream := <-r.fetchDevStreamChan:
				stream, exist := streams[hex.EncodeToString(fetchDevStream.signingKey)]
				if !exist {
					fetchDevStream.respChan <- nil
					continue
				}
				fetchDevStream.respChan <- stream
			case dApp := <-r.addDAppChan:
				dAppInstances[dApp.ID()] = dApp
			}
		}
	}()

	return r, nil

}

// start a DApp
func (r *Registry) StartDApp(dAppSigningKey ed25519.PublicKey, timeOut time.Duration) error {

	// fetch DApp
	dApp, err := r.dAppDB.Get(dAppSigningKey)
	if err != nil {
		return err
	}
	if dApp == nil {
		return fmt.Errorf("failed to fetch DApp for signing key: %x", dAppSigningKey)
	}

	// get logger
	var l *golog.Logger
	if l, err = golog.GetLogger("app name"); err != nil {
		return err
	}

	vmModules := []module.Module{
		uuidv4Mod.New(l),
		modalMod.New(l, r.api, dApp.UsedSigningKey),
		sendEthTxMod.New(r.api, l),
		randBytes.New(l),
		ethAddrMod.New(r.km),
		renderMsg.New(l),
		renderDApp.New(l),
		messageModule.New(r.msgDB, dAppSigningKey, l),
	}

	// if there is a stream for this DApp
	// we would like to mutate the logger
	// to write to the stream we have for development
	// this will write logs to the stream
	fetchDAppChan := make(chan net.Stream)
	r.fetchDevStreamChan <- fetchDAppStreamStr{
		signingKey: dAppSigningKey,
		respChan:   fetchDAppChan,
	}
	stream := <-fetchDAppChan
	if stream != nil {
		// append log module
		logger, err := loggerMod.New(stream)
		l = logger.Logger
		if err != nil {
			return err
		}
		vmModules = append(vmModules, logger)
	} else {
		if err != nil {
			return err
		}
		l.SetBackend(golog.AddModuleLevel(golog.NewLogBackend(ioutil.Discard, "", 0)))
	}

	app, err := dapp.New(l, dApp, vmModules, r.closeChan, timeOut, r.db)
	if err != nil {
		l.Error(err.Error())
		return err
	}

	// add DApp to state
	r.addDAppChan <- app

	return nil

}

func (r *Registry) fetchDApp(signingKey ed25519.PublicKey) *dapp.DApp {

	dAppRespChan := make(chan *dapp.DApp)

	// fetch DApp
	r.fetchDAppChan <- fetchDAppChanStr{
		signingKey: signingKey,
		respChan:   dAppRespChan,
	}

	return <-dAppRespChan

}

// open DApp
func (r *Registry) OpenDApp(signingKey ed25519.PublicKey, context string) error {

	dApp := r.fetchDApp(signingKey)
	if dApp == nil {
		return errors.New("it seems like that this app hasn't been started yet")
	}

	return dApp.OpenDApp(context)
}

func (r *Registry) RenderMessage(signingKey ed25519.PublicKey, payload string) (string, error) {
	dApp := r.fetchDApp(signingKey)
	if dApp == nil {
		return "", errors.New("it seems like that this app hasn't been started yet")
	}
	return dApp.RenderMessage(payload)
}

// use this to connect to a development server
func (r *Registry) ConnectDevelopmentServer(addr ma.Multiaddr) error {

	// address to peer info
	pInfo, err := pstore.InfoFromP2pAddr(addr)
	if err != nil {
		return err
	}

	// connect to peer
	if err := r.host.Connect(context.Background(), *pInfo); err != nil {
		return err
	}

	// create stream to development peer
	str, err := r.host.NewStream(context.Background(), pInfo.ID, "/dapp-development/0.0.0")
	if err != nil {
		return err
	}

	// handle stream
	r.devStreamHandler(str)

	return nil
}

// call a function in a DApp
func (r *Registry) CallFunction(signingKey ed25519.PublicKey, funcId uint, args string) error {
	dApp := r.fetchDApp(signingKey)
	if dApp == nil {
		return errors.New("it seems like that this app hasn't been started yet")
	}
	return dApp.CallFunction(funcId, args)
}

func (r *Registry) ShutDown(signingKey ed25519.PublicKey) error {
	dApp := r.fetchDApp(signingKey)
	if dApp == nil {
		return errors.New("it seems like that this app hasn't been started yet")
	}
	dApp.Close()
	return nil
}
