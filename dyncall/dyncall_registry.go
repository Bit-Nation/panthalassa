package dyncall

import "fmt"

type getCallModule struct {
	respChan chan CallModule
	id       string
}

type Registry struct {
	addModuleChan chan CallModule
	getModuleChan chan getCallModule
	closeChan     chan struct{}
}

func New() *Registry {

	r := &Registry{
		addModuleChan: make(chan CallModule),
		getModuleChan: make(chan getCallModule),
		closeChan:     make(chan struct{}),
	}

	// state
	go func() {
		modules := map[string]CallModule{}

		for {
			select {
			case <-r.closeChan:
				return
			case m := <-r.addModuleChan:
				modules[m.CallID()] = m
			case getCallMod := <-r.getModuleChan:
				mod := modules[getCallMod.id]
				getCallMod.respChan <- mod
			}
		}

	}()

	return r

}

func (r *Registry) Register(m CallModule) error {

	// exist if already registered
	respChan := make(chan CallModule)
	r.getModuleChan <- getCallModule{
		id:       m.CallID(),
		respChan: respChan,
	}
	if nil != <-respChan {
		return fmt.Errorf("a call modules with id %s has already been registered", m.CallID())
	}

	// add module
	r.addModuleChan <- m

	return nil

}

func (r *Registry) Call(callID string, payload map[string]interface{}) (map[string]interface{}, error) {

	// get call module
	respChan := make(chan CallModule)
	r.getModuleChan <- getCallModule{
		id:       callID,
		respChan: respChan,
	}
	callModule := <-respChan
	if nil == callModule {
		return map[string]interface{}{}, fmt.Errorf("a call module with call id: %s does not exist", callID)
	}

	// validate the payload
	if err := callModule.Validate(payload); err != nil {
		return map[string]interface{}{}, err
	}

	// handle the payload
	return callModule.Handle(payload)

}
