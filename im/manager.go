package im

import "sync"

var (
	globalClientManager = &ClientManager{}
)

type ClientManager struct {
	clients sync.Map
}

func (m *ClientManager) Register(client *Client) error {
	m.clients.Store(client.Key, client)
	return nil
}

func (m *ClientManager) GetClient(key string) (*Client, error) {
	client, ok := m.clients.Load(key)
	if !ok {
		return nil, nil
	}
	return client.(*Client), nil
}
