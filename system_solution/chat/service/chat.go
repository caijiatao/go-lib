package service

import (
	"golib/system_solution/chat/model"
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

func (s *chatService) PushMessage(message *model.Message) error {
	return nil
}
