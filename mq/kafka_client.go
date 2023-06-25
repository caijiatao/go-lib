package mq

import (
	"context"
	"encoding/json"
	"github.com/Shopify/sarama"
)

type kafkaClient struct {
	client sarama.Client
}

func (k *kafkaClient) Name() string {
	//TODO implement me
	panic("implement me")
}

func (k *kafkaClient) Close() error {
	//TODO implement me
	panic("implement me")
}

func (k *kafkaClient) SyncSendMessage(ctx context.Context, params interface{}) error {
	syncProducer, err := sarama.NewSyncProducerFromClient(k.client)
	if err != nil {
		return err
	}
	defer func() {
		_ = syncProducer.Close()
	}()
	ps, err := json.Marshal(params)
	if err != nil {
		return err
	}
	msg := &sarama.ProducerMessage{Topic: "test", Partition: 6, Value: sarama.StringEncoder(ps)}
	_, _, err = syncProducer.SendMessage(msg)
	if err != nil {
		return err
	}
	return nil
}

func newKafkaClient(config *Config) (*kafkaClient, error) {
	c, err := sarama.NewClient(config.Addr, config.SaramaConfig)
	if err != nil {
		return nil, err
	}
	client := &kafkaClient{
		client: c,
	}
	return client, nil
}
