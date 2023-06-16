package gorm_helper

import (
	"errors"
	"sync"
)

var (
	globalDBManager         *dbManager
	globalDBManagerInitOnce sync.Once
)

const defaultDBClientName = "default"

var (
	clientNameExists = errors.New("client name exists")
)

type dbManager struct {
	allClients sync.Map
}

func (m *dbManager) get(dbClientName string) *DB {
	db, ok := m.allClients.Load(dbClientName)
	if !ok {
		return nil
	}
	return db.(*DB)
}

func (m *dbManager) add(dbClientName string, db *DB) error {
	_, ok := m.allClients.Load(dbClientName)
	if ok {
		return clientNameExists
	}
	m.allClients.Store(dbClientName, db)
	return nil
}
