package etcd_helper

import (
	"errors"
	"sync"
)

var (
	globalManager         *etcdManager
	globalManagerInitOnce sync.Once
)

var (
	clientNameExists = errors.New("client name exists")
)

// etcdManager
// @Description: 服务所用的所有ETCD客户端管理
type etcdManager struct {
	allClients sync.Map
}

func (em *etcdManager) add(clientName string, client *EtcdClient) error {
	_, ok := em.allClients.Load(clientName)
	if ok {
		return clientNameExists
	}
	em.allClients.Store(clientName, client)
	return nil
}

func (em *etcdManager) get(clientName string) *EtcdClient {
	value, ok := em.allClients.Load(clientName)
	if !ok {
		return nil
	}
	return value.(*EtcdClient)
}

func GetDefaultClient() *EtcdClient {
	return globalManager.get(defaultClientName)
}
