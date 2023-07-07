package mq

type Message struct {
	Topic     string
	Partition int32
	Content   interface{}
}

func NewMessage(topic string, partition int32, content interface{}) *Message {
	return &Message{Topic: topic, Partition: partition, Content: content}
}
