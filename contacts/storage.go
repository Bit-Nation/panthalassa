package contacts

import (
	"github.com/asdine/storm"
)

type Contact struct {
	IdentityKey string
}

type Storage struct {
	db *storm.DB
}

func (s *Storage) All() ([]*Contact, error) {
	var contacts []*Contact
	if err := s.db.All(&contacts); err != nil {
		return nil, err
	}
	return contacts, nil
}

func (s *Storage) Save(d *Contact) error {
	return s.db.Save(d)
}

func NewStorage(db *storm.DB) *Storage {
	return &Storage{
		db: db,
	}
}
