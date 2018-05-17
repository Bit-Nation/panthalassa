package main

import (
	jsonDB "github.com/nanobox-io/golang-scribble"
)

type Store struct {
	Account Account
	DB      *jsonDB.Driver
}

func (s *Store) Send(data string) {
	logger.Info(data)
}
