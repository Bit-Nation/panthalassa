package chat

import (
	"encoding/hex"
	"encoding/json"
	"strconv"

	client "github.com/Bit-Nation/panthalassa/client"
	keyManager "github.com/Bit-Nation/panthalassa/keyManager"
	x3dh "github.com/Bit-Nation/x3dh"
	log "github.com/ipfs/go-log"
	doubleratchet "github.com/tiabc/doubleratchet"
)

var logger = log.Logger("chat")

type PangeaDoubleRachedKeyStore struct {
	dB client.PangeaKeyStoreDBInterface
	km *keyManager.KeyManager
}

func (s *PangeaDoubleRachedKeyStore) Get(k doubleratchet.Key, msgNum uint) (mk doubleratchet.Key, ok bool) {

	msgNumStr := strconv.FormatUint(uint64(msgNum), 10)
	keyStr := hex.EncodeToString(k[:])

	msgKeyCipher := s.dB.Get(keyStr, msgNumStr)

	if msgKeyCipher == "" {
		return doubleratchet.Key{}, false
	}

	msgKeyStr, err := s.km.AESDecrypt(msgKeyCipher)
	if err != nil {
		logger.Fatal("Failed to decrypt message key")
		return doubleratchet.Key{}, false
	}

	rawMsgKey, err := hex.DecodeString(msgKeyStr)
	if err != nil {
		return doubleratchet.Key{}, false
	}

	msgKey := doubleratchet.Key{}
	copy(msgKey[:], rawMsgKey[:32])

	return msgKey, true

}

func (s *PangeaDoubleRachedKeyStore) Put(k doubleratchet.Key, msgNum uint, mk doubleratchet.Key) {

	keyStr := hex.EncodeToString(k[:])
	msgNumStr := strconv.FormatUint(uint64(msgNum), 10)
	messageKey := hex.EncodeToString(mk[:])

	s.dB.Put(keyStr, msgNumStr, messageKey)

}

func (s *PangeaDoubleRachedKeyStore) DeleteMk(k doubleratchet.Key, msgNum uint) {

	keyStr := hex.EncodeToString(k[:])
	msgNumStr := strconv.FormatUint(uint64(msgNum), 10)

	s.dB.DeleteMk(keyStr, msgNumStr)

}

func (s *PangeaDoubleRachedKeyStore) DeletePk(k doubleratchet.Key) {

	keyStr := hex.EncodeToString(k[:])
	s.dB.DeletePk(keyStr)

}

func (s *PangeaDoubleRachedKeyStore) Count(k doubleratchet.Key) uint {

	keyStr := hex.EncodeToString(k[:])

	countStr := s.dB.Count(keyStr)

	// @todo skipping the error is not so nice but we can't handle it duo the limits of the interface
	c, err := strconv.ParseUint(countStr, 10, 32)
	if err != nil {
		logger.Error("Failed to parse int: ", err)
	}

	return uint(c)
}

func (s *PangeaDoubleRachedKeyStore) All() map[doubleratchet.Key]map[uint]doubleratchet.Key {

	m := map[doubleratchet.Key]map[uint]doubleratchet.Key{}

	all := s.dB.All()

	if err := json.Unmarshal([]byte(all), &m); err != nil {
		// @todo handle error better
		logger.Error("Failed to unmarshal all: ", err)
		return m
	}

	return m

}

type PangeaOneTimePreKeyStore struct {
	dB client.OneTimePreKeyStoreDBInterface
	km keyManager.KeyManager
}

// put a key to the store
func (s *PangeaOneTimePreKeyStore) Put(keyPair x3dh.KeyPair) error {

	pub := hex.EncodeToString(keyPair.PublicKey[:])
	privStr := hex.EncodeToString(keyPair.PrivateKey[:])

	cipherText, err := s.km.AESEncrypt(privStr)
	if err != nil {
		return err
	}

	return s.dB.Put(pub, cipherText)

}

func (s *PangeaOneTimePreKeyStore) Get(key x3dh.PublicKey) (x3dh.PrivateKey, error) {

	pub := hex.EncodeToString(key[:])

	privCipher, err := s.dB.Get(pub)
	if err != nil {
		return x3dh.PrivateKey{}, err
	}

	plainPriv, err := s.km.AESDecrypt(privCipher)
	if err != nil {
		return x3dh.PrivateKey{}, err
	}

	rawPriv, err := hex.DecodeString(plainPriv)
	if err != nil {
		return x3dh.PrivateKey{}, err
	}

	var priv [32]byte
	copy(priv[:], rawPriv[:32])

	return priv, nil

}

func (s *PangeaOneTimePreKeyStore) Has(key x3dh.PublicKey) (bool, error) {
	pub := hex.EncodeToString(key[:])
	return s.dB.Has(pub)
}

func (s *PangeaOneTimePreKeyStore) Delete(key x3dh.PublicKey) error {
	pub := hex.EncodeToString(key[:])
	return s.dB.Delete(pub)
}
