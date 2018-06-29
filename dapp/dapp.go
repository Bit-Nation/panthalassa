package dapp

import (
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	module "github.com/Bit-Nation/panthalassa/dapp/module"
	cbModule "github.com/Bit-Nation/panthalassa/dapp/module/callbacks"
	dAppRenderer "github.com/Bit-Nation/panthalassa/dapp/module/renderer/dapp"
	msgRenderer "github.com/Bit-Nation/panthalassa/dapp/module/renderer/message"
	logger "github.com/op/go-logging"
	otto "github.com/robertkrimen/otto"
)

type DApp struct {
	vm           *otto.Otto
	logger       *logger.Logger
	app          *JsonRepresentation
	closeChan    chan<- *JsonRepresentation
	dAppRenderer *dAppRenderer.Module
	msgRenderer  *msgRenderer.Module
	cbMod        *cbModule.Module
}

// close DApp
func (d *DApp) Close() {
	d.vm.Interrupt <- func() {
		d.logger.Info(fmt.Sprintf("shutting down: %s (%s)", hex.EncodeToString(d.app.SignaturePublicKey), d.app.Name))
		d.closeChan <- d.app
	}
}

func (d *DApp) ID() string {
	return hex.EncodeToString(d.app.SignaturePublicKey)
}

func (d *DApp) RenderDApp(context string) error {
	return d.dAppRenderer.OpenDApp(context)
}

func (d *DApp) RenderMessage(payload string) (string, error) {
	return d.msgRenderer.RenderMessage(payload)
}

func (d *DApp) CallFunction(id uint, args string) error {
	return d.cbMod.CallFunction(id, args)
}

// will start a DApp based on the given config file
//
func New(l *logger.Logger, app *JsonRepresentation, vmModules []module.Module, closer chan<- *JsonRepresentation, timeOut time.Duration) (*DApp, error) {

	// check if app is valid
	valid, err := app.VerifySignature()
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, InvalidSignature
	}

	// create VM
	vm := otto.New()
	vm.Interrupt = make(chan func(), 1)

	// register all vm modules
	for _, m := range vmModules {
		if err := m.Register(vm); err != nil {
			return nil, err
		}
	}

	// register DApp renderer
	dr := dAppRenderer.New(l)
	if err := dr.Register(vm); err != nil {
		return nil, err
	}

	// register message renderer
	mr := msgRenderer.New(l)
	if err := dr.Register(vm); err != nil {
		return nil, err
	}

	// register callbacks module
	cbm := cbModule.New(l)
	if err := cbm.Register(vm); err != nil {
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
	}

	wait := make(chan error, 1)

	// start the DApp async
	go func() {
		_, err := vm.Run(app.Code)
		fmt.Println(err)
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
