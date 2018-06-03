package panthalassa

// those interfaces are copies to make them available during gomobile

type PangeaKeyStoreDBInterface interface {
	// get the encryption secret by a key and a message
	// string must be "" if no key is found for the message
	Get(key string, msgNum string) string

	// save a messageKey in a mapping between msgNum under
	// key => mapping(msgNum => messageKey)
	Put(key string, msgNum string, messageKey string)

	// delete a messageKey by a key and msgNum
	DeleteMk(key string, msgNum string)

	// delete all messages under a key
	DeletePk(key string)

	// count all messages of key
	Count(key string) string

	// fetch all message keys
	// should be returned as json
	// {
	//		"key_one" => {
	// 			1 => "encrypted_key_blabla"
	// 		},
	// 		"key_two" => {
	// 			1 => "encrypted_key_jngjfj"
	//		}
	// }
	All() string
}

type OneTimePreKeyStoreDBInterface interface {
	Put(pubKey string, encPrivKey string) error
	Get(pubKey string) (string, error)
	Has(pubKey string) (bool, error)
	Delete(pubKey string) error
}
