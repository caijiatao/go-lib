package mq

import "github.com/Shopify/sarama"

type Config struct {
	ClientName   string
	Addr         []string
	SaramaConfig *sarama.Config

	RetryCount int
}

func NewConfig(clientName string, addr []string) *Config {
	config := &Config{
		ClientName:   clientName,
		Addr:         addr,
		SaramaConfig: sarama.NewConfig(),
		RetryCount:   3,
	}
	return config
}

func readConfig() []*Config {
	return nil
}
