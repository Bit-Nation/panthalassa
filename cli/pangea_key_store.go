package main

import (
	"encoding/json"
	"strconv"
)

func NewKeyStore() PangeaKeyStoreDB {
	return PangeaKeyStoreDB{
		db: make(map[string]map[string]string),
	}
}

type PangeaKeyStoreDB struct {
	db map[string]map[string]string
}

func (p *PangeaKeyStoreDB) Get(key, msgNum string) string {
	return p.db[key][msgNum]
}

func (p *PangeaKeyStoreDB) Put(key string, msgNum string, messageKey string) {
	p.db[key][msgNum] = messageKey
}

func (p *PangeaKeyStoreDB) DeleteMk(key string, msgNum string) {
	delete(p.db[key], msgNum)
}

func (p *PangeaKeyStoreDB) DeletePk(key string) {
	delete(p.db, key)
}

func (p *PangeaKeyStoreDB) Count(key string) string {
	return strconv.Itoa(len(p.db[key]))
}

func (p *PangeaKeyStoreDB) All() string {
	data, err := json.Marshal(p.db)
	if err != nil {
		panic(err)
	}
	return string(data)
}
