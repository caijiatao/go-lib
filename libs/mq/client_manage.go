package mq

import (
	"errors"
	"sync"
)

var (
	globalClientManage         *clientManage
	globalClientManageInitOnce sync.Once
	clientNameExists           = errors.New("client name exists")
)

type clientManage struct {
	allClients sync.Map
}

func (c *clientManage) register(client Client) error {
	_, ok := c.allClients.Load(client.Name())
	if ok {
		return clientNameExists
	}
	c.allClients.Store(client.Name(), client)
	return nil
}

func (c *clientManage) get(name string) Client {
	client, ok := c.allClients.Load(name)
	if ok {
		return client.(Client)
	}
	return nil
}

func (c *clientManage) Close() (err error) {
	c.allClients.Range(func(clientName, v any) bool {
		client := v.(Client)
		_ = client.Close()
		c.allClients.Delete(client.Name())
		return true
	})
	return nil
}
