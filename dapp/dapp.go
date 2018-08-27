package dapp

import (
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	module "github.com/Bit-Nation/panthalassa/dapp/module"
	cbModule "github.com/Bit-Nation/panthalassa/dapp/module/callbacks"
	dbModule "github.com/Bit-Nation/panthalassa/dapp/module/db"
	dAppRenderer "github.com/Bit-Nation/panthalassa/dapp/module/renderer/dapp"
	msgRenderer "github.com/Bit-Nation/panthalassa/dapp/module/renderer/message"
	bolt "github.com/coreos/bbolt"
	log "github.com/ipfs/go-log"
	logger "github.com/op/go-logging"
	duktape "gopkg.in/olebedev/go-duktape.v3"
)

var sysLog = log.Logger("dapp")

type DApp struct {
	vm     *duktape.Context
	logger *logger.Logger
	app    *Data
	// will be called when the app shut down
	closeChan    chan<- *Data
	dAppRenderer *dAppRenderer.Module
	msgRenderer  *msgRenderer.Module
	cbMod        *cbModule.Module
	dbMod        *dbModule.BoltStorage
	vmModules    []module.Module
}

// close DApp
func (d *DApp) Close() {
	d.vm.Interrupt <- func() {
		d.logger.Info(fmt.Sprintf("shutting down: %s (%s)", hex.EncodeToString(d.app.UsedSigningKey), d.app.Name))
		for _, mod := range d.vmModules {
			if err := mod.Close(); err != nil {
				sysLog.Error(err)
			}
		}
		d.closeChan <- d.app
	}
}

func (d *DApp) ID() string {
	return hex.EncodeToString(d.app.UsedSigningKey)
}

func (d *DApp) OpenDApp(context string) error {
	return d.dAppRenderer.OpenDApp(context)
}

func (d *DApp) RenderMessage(payload string) (string, error) {
	return d.msgRenderer.RenderMessage(payload)
}

func (d *DApp) CallFunction(id uint, args string) error {
	return d.cbMod.CallFunction(id, args)
}

// will start a DApp based on the given config file
func New(l *logger.Logger, app *Data, vmModules []module.Module, closer chan<- *Data, timeOut time.Duration, db *bolt.DB) (*DApp, error) {

	// check if app is valid
	valid, err := app.VerifySignature()
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, InvalidSignature
	}

	// create VM
	vm := duktape.New()
	vm.Interrupt = make(chan func(), 1)

	// register all vm modules
	for _, m := range vmModules {
		if err := m.Register(vm); err != nil {
			return nil, err
		}
	}

	// register DApp renderer
	dr := dAppRenderer.New(l)
	vmModules = append(vmModules, dr)
	if err := dr.Register(vm); err != nil {
		return nil, err
	}

	// register message renderer
	mr := msgRenderer.New(l)
	vmModules = append(vmModules, mr)
	if err := mr.Register(vm); err != nil {
		return nil, err
	}

	// register callbacks module
	cbm := cbModule.New(l)
	vmModules = append(vmModules, cbm)
	if err := cbm.Register(vm); err != nil {
		return nil, err
	}

	// register database (the one used by the dapp)
	dAppDBStorage, err := dbModule.NewBoltStorage(db, app.UsedSigningKey)
	if err != nil {
		return nil, err
	}
	dbm := dbModule.New(dAppDBStorage, l)
	vmModules = append(vmModules, dbm)
	if err := dbm.Register(vm); err != nil {
		return nil, err
	}

	dApp := &DApp{
		vm:           vm,
		logger:       l,
		app:          app,
		closeChan:    closer,
		dAppRenderer: dr,
		msgRenderer:  mr,
		cbMod:        cbm,
		dbMod:        dAppDBStorage,
		vmModules:    vmModules,
	}

	wait := make(chan error, 1)

	// start the DApp async
	go func() {
		_, err := vm.Run(app.Code)
		if err != nil {
			l.Errorf(err.Error())
		}
		wait <- err
	}()

	// wait for the DApp with given timeout
	select {
	case err := <-wait:
		if err != nil {
			return nil, err
		}
		return dApp, nil
	case <-time.After(timeOut):
		vm.Interrupt <- func() {}
		closer <- app
		return nil, errors.New("timeout - failed to start DApp")
	}

}
