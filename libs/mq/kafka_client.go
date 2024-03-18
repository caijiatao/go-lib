package mq

import (
	"context"
	"encoding/json"
	"github.com/Shopify/sarama"
	"golib/libs/logger"
)

type KafkaClient struct {
	config *Config
	client sarama.Client
}

func (k *KafkaClient) Name() string {
	return k.config.ClientName
}

func (k *KafkaClient) Close() error {
	err := k.client.Close()
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (k *KafkaClient) SyncSendMessage(ctx context.Context, message *Message) error {
	syncProducer, err := sarama.NewSyncProducerFromClient(k.client)
	if err != nil {
		return err
	}
	defer func() {
		_ = syncProducer.Close()
	}()
	ps, err := json.Marshal(message.Content)
	if err != nil {
		return err
	}
	msg := &sarama.ProducerMessage{Topic: message.Topic, Partition: message.Partition, Value: sarama.StringEncoder(ps)}
	for i := 0; i < k.config.RetryCount; i++ {
		_, _, err = syncProducer.SendMessage(msg)
		if err != nil {
			logger.Error(err)
			continue
		}
	}
	if err != nil {
		return err
	}
	return nil
}

func NewKafkaClient(config *Config) (*KafkaClient, error) {
	c, err := sarama.NewClient(config.Addr, config.SaramaConfig)
	if err != nil {
		return nil, err
	}
	client := &KafkaClient{
		config: config,
		client: c,
	}
	return client, nil
}
