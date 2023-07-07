package mq

import (
	"context"
)

// Client
// @Description: 依赖client interface ，后续可以直接替换MQ底层client ，目前默认为kafka
type Client interface {
	Name() string
	SyncSendMessage(ctx context.Context, params *Message) error
	Close() error
}

func Init() (err error) {
	configs := readConfig()
	defer func() {
		if err != nil {
			_ = Close()
		}
	}()
	for _, config := range configs {
		_, err = NewClient(config)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewClient(config *Config) (Client, error) {
	globalClientManageInitOnce.Do(func() {
		globalClientManage = &clientManage{}
	})
	// 默认创建kafka的client
	client, err := NewKafkaClient(config)
	if err != nil {
		return nil, err
	}
	err = globalClientManage.register(client)
	if err != nil {
		return client, err
	}
	return client, nil
}

func GetClient(clientName string) Client {
	return globalClientManage.get(clientName)
}

func Close() error {
	if globalClientManage == nil {
		return nil
	}
	err := globalClientManage.Close()
	return err
}
