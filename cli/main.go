package main

import (
	"os"

	km "github.com/Bit-Nation/panthalassa/keyManager"
	ks "github.com/Bit-Nation/panthalassa/keyStore"
	mnemonic "github.com/Bit-Nation/panthalassa/mnemonic"
	panthalassa "github.com/Bit-Nation/panthalassa/mobile"
	log "github.com/ipfs/go-log"
	jsonDB "github.com/nanobox-io/golang-scribble"
	uuid "github.com/satori/go.uuid"
	iShell "gopkg.in/abiosoft/ishell.v2"
)

const DevRendezvousKey = "akhgp58izorhalsdipfo3uh5orpawoudshfalskduf43topa"
const LogFile = "log.out"
const DBName = ".database"

var logger = log.Logger("hi")

type Account struct {
	Name         string `json:"name"`
	AccountStore string `json:"account_store"`
}

func main() {

	//Database
	db, err := jsonDB.New(DBName, nil)
	if err != nil {
		panic(err)
	}

	shell := iShell.New()

	// register command to start panthalassa
	shell.AddCmd(&iShell.Cmd{
		Name: "start",
		Help: "start panthalassa",
		Func: func(c *iShell.Context) {

			mne, err := mnemonic.New()
			if err != nil {
				c.Err(err)
				return
			}

			store, err := ks.NewFromMnemonic(mne)
			if err != nil {
				c.Err(err)
				return
			}

			keyManager := km.CreateFromKeyStore(store)
			keyStoreStr, err := keyManager.Export("pw", "pw")

			err = panthalassa.Start(keyStoreStr, "pw", DevRendezvousKey, nil)
			if err != nil {
				c.Err(err)
				return
			}

			c.Println("Started panthalassa")

		},
	})

	// stop panthalassa
	shell.AddCmd(&iShell.Cmd{
		Name: "stop",
		Help: "stop's the current panthalassa instance",
		Func: func(c *iShell.Context) {
			err := panthalassa.Stop()
			if err != nil {
				c.Err(err)
				return
			}
			c.Println("stopped panthalassa")
		},
	})

	// display private key
	shell.AddCmd(&iShell.Cmd{
		Name: "eth:private",
		Help: "show's the ethereum private key",
		Func: func(c *iShell.Context) {
			pk, err := panthalassa.EthPrivateKey()
			if err != nil {
				c.Err(err)
				return
			}
			c.Println("your private key is: ", pk)
		},
	})

	// display address
	shell.AddCmd(&iShell.Cmd{
		Name: "eth:address",
		Help: "display ethereum address",
		Func: func(c *iShell.Context) {
			addr, err := panthalassa.EthAddress()
			if err != nil {
				c.Err(err)
				return
			}
			logger.Error("hi")

			c.Println("your ethereum address is:", addr)
		},
	})

	shell.AddCmd(&iShell.Cmd{
		Name: "log:debug",
		Help: "Enable debug logging",
		Func: func(c *iShell.Context) {
			f, err := os.Create(LogFile)
			if err != nil {
				c.Err(err)
				return
			}
			log.Configure(log.Output(f), log.LevelDebug)
			c.Println("Enabled logging (for debug). Output file: ", f.Name())
		},
	})

	shell.AddCmd(&iShell.Cmd{
		Name: "log:warn",
		Help: "Enable debug logging",
		Func: func(c *iShell.Context) {
			f, err := os.Create(LogFile)
			if err != nil {
				c.Err(err)
				return
			}
			log.Configure(log.Output(f))
			//2 = WARNING
			log.SetAllLoggers(2)
			c.Println("Enabled logging (for warning's). Output file: ", f.Name())
		},
	})

	shell.AddCmd(&iShell.Cmd{
		Name: "account:new",
		Help: "Create a new Account",
		Func: func(c *iShell.Context) {

			//get username
			c.Println("Account name: ")
			accountName := c.ReadLine()

			// get password
			c.Print("Password for account: ")
			password := c.ReadLine()

			// create mnemonic
			mne, err := mnemonic.New()
			if err != nil {
				c.Err(err)
				return
			}

			// create key store form mnemonic
			store, err := ks.NewFromMnemonic(mne)
			if err != nil {
				c.Err(err)
				return
			}

			// create key manager from key store
			keyManager := km.CreateFromKeyStore(store)
			exportedAccount, err := keyManager.Export(password, password)
			if err != nil {
				c.Err(err)
				return
			}

			// uuid
			id, err := uuid.NewV4()
			if err != nil {
				c.Err(err)
				return
			}

			err = db.Write("account", id.String(), &Account{
				Name:         accountName,
				AccountStore: exportedAccount,
			})

			if err != nil {
				c.Err(err)
				return
			}

			c.Println("safed account store")

		},
	})

	// run shell
	shell.Run()
}
