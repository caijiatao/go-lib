package orm

import (
	"errors"
	"sync"
)

var (
	globalClientManager = &clientManager{}
)

const defaultDBClientName = "default"

var (
	clientNameExists = errors.New("client name exists")
)

type clientManager struct {
	allClients sync.Map
}

func (m *clientManager) get(dbClientName string) Client {
	db, ok := m.allClients.Load(dbClientName)
	if !ok {
		return nil
	}
	return db.(*GormClientProxy)
}

func (m *clientManager) add(dbClientName string, db Client) error {
	_, ok := m.allClients.Load(dbClientName)
	if ok {
		return clientNameExists
	}
	m.allClients.Store(dbClientName, db)
	return nil
}
