package db

import (
	"encoding/binary"

	aes "github.com/Bit-Nation/panthalassa/crypto/aes"
	km "github.com/Bit-Nation/panthalassa/keyManager"
	storm "github.com/asdine/storm"
	sq "github.com/asdine/storm/q"
	log "github.com/ipfs/go-log"
	dr "github.com/tiabc/doubleratchet"
)

type DRKey struct {
	ID     int            `storm:"id,increment"`
	Key    dr.Key         `storm:"index"`
	MsgNum uint           `storm:"index"`
	MsgKey aes.CipherText `storm:"index"`
}

type DRKeyStorage interface {
	Get(k dr.Key, msgNum uint) (mk dr.Key, ok bool)
	Put(k dr.Key, msgNum uint, mk dr.Key)
	DeleteMk(k dr.Key, msgNum uint)
	DeletePk(k dr.Key)
	Count(k dr.Key) uint
	All() map[dr.Key]map[uint]dr.Key
}

type BoltDRKeyStorage struct {
	km *km.KeyManager
	db storm.Node
}

func NewBoltDRKeyStorage(db storm.Node, km *km.KeyManager) *BoltDRKeyStorage {
	return &BoltDRKeyStorage{
		km: km,
		db: db,
	}
}

var logger = log.Logger("database")

func uintToBytes(uint uint) []byte {
	num := make([]byte, 8)
	binary.LittleEndian.PutUint64(num, uint64(uint))
	return num
}

func bytesToUint(uint []byte) uint64 {
	return binary.LittleEndian.Uint64(uint)
}

func (s *BoltDRKeyStorage) Get(k dr.Key, msgNum uint) (mk dr.Key, ok bool) {

	// fetch dr key
	q := s.db.Select(sq.And(
		sq.Eq("Key", k),
		sq.Eq("MsgNum", msgNum),
	))

	// count
	amount, err := q.Count(&DRKey{})
	if amount <= 0 {
		return dr.Key{}, false
	}

	drKey := &DRKey{}
	if err := q.First(drKey); err != nil {
		logger.Error(err)
		return dr.Key{}, false
	}

	mk = dr.Key{}
	plainMK, err := s.km.AESDecrypt(drKey.MsgKey)
	if err != nil {
		logger.Error(err)
		return dr.Key{}, false
	}
	if err := q.First(drKey); err != nil {
		logger.Error(err)
		return dr.Key{}, false
	}
	if len(plainMK) != 32 {
		logger.Errorf("got invalid message key with len != 32")
		return dr.Key{}, false
	}
	copy(mk[:], plainMK)

	return mk, true

}

func (s *BoltDRKeyStorage) Put(k dr.Key, msgNum uint, mk dr.Key) {

	ct, err := s.km.AESEncrypt(mk[:])
	if err != nil {
		logger.Error(err)
		return
	}

	// persist double ratchet key
	err = s.db.Save(&DRKey{
		Key:    k,
		MsgNum: msgNum,
		MsgKey: ct,
	})

	// @todo we need to change the dr package to use a better interface
	if err != nil {
		logger.Error(err)
	}

}

// @todo we need to change the dr package to use a better interface (error handling)
func (s *BoltDRKeyStorage) DeleteMk(k dr.Key, msgNum uint) {

	// fetch double ratchet key
	q := s.db.Select(sq.And(
		sq.Eq("Key", k),
		sq.Eq("MsgNum", msgNum),
	))
	var drKey DRKey
	if err := q.First(&drKey); err != nil {
		logger.Error(err)
		return
	}

	// delete it
	if err := s.db.DeleteStruct(&drKey); err != nil {
		logger.Error(err)
		return
	}

}

func (s *BoltDRKeyStorage) DeletePk(k dr.Key) {

	// delete all DR keys under the given key
	q := s.db.Select(sq.Eq("Key", k))

	var drKey DRKey
	// @todo we need to change the dr package to use a better interface
	if err := q.Delete(&drKey); err != nil {
		logger.Error(err)
	}

}

func (s *BoltDRKeyStorage) Count(k dr.Key) uint {

	q := s.db.Select(sq.Eq("Key", k))
	count, err := q.Count(&DRKey{})
	if err != nil {
		logger.Error(err)
		return 0
	}

	return uint(count)

}

func (s *BoltDRKeyStorage) All() map[dr.Key]map[uint]dr.Key {

	keys := map[dr.Key]map[uint]dr.Key{}

	allDRKeys := []DRKey{}
	if err := s.db.All(&allDRKeys); err != nil {
		logger.Error(err)
		return map[dr.Key]map[uint]dr.Key{}
	}

	for _, drKey := range allDRKeys {
		drk := dr.Key{}
		plainDRKey, err := s.km.AESDecrypt(drKey.MsgKey)
		if err != nil {
			logger.Error(err)
			return map[dr.Key]map[uint]dr.Key{}
		}
		if len(plainDRKey) != 32 {
			logger.Error("got invalid plain key (invalid len)")
			return map[dr.Key]map[uint]dr.Key{}
		}
		copy(drk[:], plainDRKey)

		if _, exist := keys[drKey.Key]; !exist {
			keys[drKey.Key] = map[uint]dr.Key{}
		}

		keys[drKey.Key][drKey.MsgNum] = drk
	}

	return keys

}
