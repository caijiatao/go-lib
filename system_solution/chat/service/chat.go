package service

import (
	"sync"
)

var (
	chatServiceInstance *chatService
	chatServiceOnce     sync.Once
)

type chatService struct{}

func ChatService() *chatService {
	chatServiceOnce.Do(func() {
		chatServiceInstance = &chatService{}
	})
	return chatServiceInstance
}

func (s *chatService) PushMessage(message *Message) error {
	return nil
}
