package main

import (
	"os"

	km "github.com/Bit-Nation/panthalassa/keyManager"
	ks "github.com/Bit-Nation/panthalassa/keyStore"
	mnemonic "github.com/Bit-Nation/panthalassa/mnemonic"
	panthalassa "github.com/Bit-Nation/panthalassa/mobile"
	log "github.com/ipfs/go-log"
	ishell "gopkg.in/abiosoft/ishell.v2"
)

const DevRendezvousKey = "akhgp58izorhalsdipfo3uh5orpawoudshfalskduf43topa"
const LogFile = "log.out"

var logger = log.Logger("hi")

func main() {

	shell := ishell.New()

	// display welcome info.
	shell.Println("Panthalassa interactive shell")

	// register command to start panthalassa
	shell.AddCmd(&ishell.Cmd{
		Name: "start",
		Help: "start panthalassa",
		Func: func(c *ishell.Context) {

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
	shell.AddCmd(&ishell.Cmd{
		Name: "stop",
		Help: "stop's the current panthalassa instance",
		Func: func(c *ishell.Context) {
			err := panthalassa.Stop()
			if err != nil {
				c.Err(err)
				return
			}
			c.Println("stopped panthalassa")
		},
	})

	// display private key
	shell.AddCmd(&ishell.Cmd{
		Name: "eth:private",
		Help: "show's the ethereum private key",
		Func: func(c *ishell.Context) {
			pk, err := panthalassa.EthPrivateKey()
			if err != nil {
				c.Err(err)
				return
			}
			c.Println("your private key is: ", pk)
		},
	})

	// display address
	shell.AddCmd(&ishell.Cmd{
		Name: "eth:address",
		Help: "display ethereum address",
		Func: func(c *ishell.Context) {
			addr, err := panthalassa.EthAddress()
			if err != nil {
				c.Err(err)
				return
			}
			logger.Error("hi")

			c.Println("your ethereum address is:", addr)
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "log:debug",
		Help: "Enable debug logging",
		Func: func(c *ishell.Context) {
			f, err := os.Create(LogFile)
			if err != nil {
				c.Err(err)
				return
			}
			log.Configure(log.Output(f), log.LevelDebug)
			c.Println("Enabled logging (for debug). Output file: ", f.Name())
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "log:warn",
		Help: "Enable debug logging",
		Func: func(c *ishell.Context) {
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

	// run shell
	shell.Run()
}
