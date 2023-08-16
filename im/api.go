package im

import "golib/logger"

func PushMessage(message *Message) {
	client, err := globalClientManager.GetClient(message.Key)
	if err != nil {
		logger.Errorf("get client error: %v", err)
		return
	}
	if client == nil {
		logger.Warnf("client is nil")
		return
	}
	client.send <- message.Msg
}
