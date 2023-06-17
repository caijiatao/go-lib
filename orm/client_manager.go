package orm

import (
	"errors"
	"sync"
)

var (
	globalDBManager         *ormManger
	globalDBManagerInitOnce sync.Once
)

const defaultDBClientName = "default"

var (
	clientNameExists = errors.New("client name exists")
)

type ormManger struct {
	allClients sync.Map
}

func (m *ormManger) get(dbClientName string) Client {
	db, ok := m.allClients.Load(dbClientName)
	if !ok {
		return nil
	}
	return db.(Client)
}

func (m *ormManger) add(clientName string, client Client) error {
	_, ok := m.allClients.Load(clientName)
	if ok {
		return clientNameExists
	}
	m.allClients.Store(clientName, client)
	return nil
}
