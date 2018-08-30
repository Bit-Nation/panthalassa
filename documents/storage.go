package documents

import (
	"github.com/asdine/storm"
	"github.com/Bit-Nation/panthalassa/crypto/aes"
	"github.com/Bit-Nation/panthalassa/keyManager"
)

type Document struct {
	ID int `storm:"id,increment"`
	Title string
	MimeType string
	Content []byte `json:"-"`
	EncryptedContent aes.CipherText
	CreatedAt int64
}

type Storage struct {
	db *storm.DB
	km *keyManager.KeyManager
}

func (s *Storage) All() ([]*Document, error) {
	var docs []*Document
	if err := s.db.All(&docs); err != nil {
		return nil, err
	}
	return docs, nil
}

func (s *Storage) Save(d *Document) error {
	// encrypt content
	ct, err := s.km.AESEncrypt(d.Content)
	if err != nil {
		return err
	}
	d.EncryptedContent = ct
	
	return s.db.Save(d)
}

func NewStorage(db *storm.DB, km *keyManager.KeyManager) *Storage {
	return &Storage{
		db: db,
		km: km,
	}
}

