package mq

import "github.com/Shopify/sarama"

type Config struct {
	ClientName   string
	Addr         []string
	SaramaConfig *sarama.Config
}
