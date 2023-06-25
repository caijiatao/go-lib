package mq

import "context"

// Client
// @Description: 依赖client interface ，后续可以直接替换MQ底层client ，目前默认为kafka
type Client interface {
	Name() string
	SyncSendMessage(ctx context.Context, params interface{}) error
	Close() error
}

func NewClient(config *Config) (Client, error) {
	globalClientManageInitOnce.Do(func() {
		globalClientManage = &clientManage{}
	})
	// 默认创建kafka的client
	client, err := newKafkaClient(config)
	if err != nil {
		return nil, err
	}
	err = globalClientManage.register(client)
	if err != nil {
		return client, err
	}
	return client, nil
}

func Close() error {
	err := globalClientManage.Close()
	return err
}
