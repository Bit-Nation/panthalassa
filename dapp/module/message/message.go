package message

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"

	reqLim "github.com/Bit-Nation/panthalassa/dapp/request_limitation"
	validator "github.com/Bit-Nation/panthalassa/dapp/validator"
	db "github.com/Bit-Nation/panthalassa/db"
	log "github.com/ipfs/go-log"
	logger "github.com/op/go-logging"
	ed25519 "golang.org/x/crypto/ed25519"
	duktape "gopkg.in/olebedev/go-duktape.v3"
)

type Module struct {
	msgStorage db.ChatMessageStorage
	dAppPubKey ed25519.PublicKey
	logger     *logger.Logger
	reqLim     *reqLim.CountThrottling
}

var sysLog = log.Logger("messsage")

func New(msgStorage db.ChatMessageStorage, dAppPubKey ed25519.PublicKey, logger *logger.Logger) *Module {
	return &Module{
		msgStorage: msgStorage,
		dAppPubKey: dAppPubKey,
		logger:     logger,
		reqLim:     reqLim.NewCountThrottling(4, time.Second*60, 10, errors.New("send message queue is full")),
	}
}

func hasKey(stack []string, search string) bool {
	for _, v := range stack {
		if v == search {
			return true
		}
	}
	return false
}

func (m *Module) Close() error {
	return nil
}

func (m *Module) Register(vm *duktape.Context) error {
	_, err := vm.PushGlobalGoFunction("sendMessage", func(context *duktape.Context) int {
		var itemsToPopBeforeCallback int
		sysLog.Debug("send message")
		// validate function call
		v := validator.New()
		// first param must be the chat
		v.Set(0, &validator.TypeString)
		// second param must be the message payload
		v.Set(1, &validator.TypeObject)
		// callback
		v.Set(2, &validator.TypeFunction)
		// utils to handle an occurred error
		handleError := func(errMsg string) int {
			if context.IsFunction(2) {
				context.PopN(itemsToPopBeforeCallback)
				context.PushString(errMsg)
				context.Call(1)
				return 0
			}
			m.logger.Error(errMsg)
			return 1
		}
		if err := v.Validate(vm); err != nil {
			// in the case an callback has been passed we want to call it with the error
			return handleError(err.Error())
		}
		objv := validator.NewObjValidator()
		objv.Set("shouldSend", validator.ObjTypeBool, true)
		objv.Set("params", validator.ObjTypeObject, false)
		objv.Set("type", validator.ObjTypeString, false)
		if err := objv.Validate(vm, 1); err != nil {
			return handleError(err.Error())
		}

		dAppMessage := db.DAppMessage{
			DAppPublicKey: m.dAppPubKey,
			Params:        map[string]interface{}{},
		}

		// chat in which the message should be persisted
		chatStr := context.ToString(0)
		itemsToPopBeforeCallback++
		chat, err := hex.DecodeString(chatStr)
		if err != nil {
			return handleError(err.Error())
		}
		if len(chat) != 32 {
			return handleError("chat must be 32 bytes long")
		}
		// should this message be sent or kept locally?
		if !context.GetPropString(1, "shouldSend") {
			handleError(`key "shouldSend" doesn't exist`)
		}
		dAppMessage.ShouldSend = context.ToBoolean(-1)
		itemsToPopBeforeCallback++
		itemsToPopBeforeCallback++
		// set optional type of the message
		if !context.GetPropString(1, "type") {
			handleError(`key "type" doesn't exist`)
		}
		dAppMessage.Type = context.ToString(-1)
		itemsToPopBeforeCallback++
		itemsToPopBeforeCallback++
		// set optional params
		if !context.GetPropString(1, "params") {
			handleError(`key "params" doesn't exist`)
		}

		//@TODO Find a way to iterate over nested object key / values
		//@TODO Lets pass objects to VM in json format and handle them fully in golang as it's much more powerful?
		//@TODO Figure out if Duktape has a way to iterate over object keys and values
		if !context.GetPropString(-1, "key") {
			handleError(`key "key" doesn't exist`)
		}

		vmValue := context.ToString(-1)
		itemsToPopBeforeCallback++
		itemsToPopBeforeCallback++
		dAppMessage.Params["key"] = vmValue
		// marshal params
		marshaledParams, err := json.Marshal(dAppMessage.Params)
		if err != nil {
			return handleError(err.Error())
		}

		// make sure it's less than 64 KB
		if len(marshaledParams) > 64*1024 {
			return handleError("the message params can't be bigger than 64 kb")
		}
		if err := json.Unmarshal(marshaledParams, &dAppMessage.Params); err != nil {
			return handleError(err.Error())
		}

		throttlingFunc := func(dec chan struct{}) {
			defer func() {
				dec <- struct{}{}
			}()
			// persist message
			if err := m.msgStorage.PersistDAppMessage(chat, dAppMessage); err != nil {
				handleError(err.Error())
			}
			//@TODO FIND OUT WHY DO WE NEED TO CLEAR THE STACK LIKE THIS?
			context.PopN(itemsToPopBeforeCallback)
			context.PushUndefined()
			context.Call(1)
			return
		}
		dec := make(chan struct{}, 1)
		throttlingFunc(dec)
		//@TODO Find a way to fix throttling
		//m.reqLim.Exec(throttlingFunc)
		return 0
	})
	return err
}
