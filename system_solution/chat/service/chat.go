package service

import (
	"context"
	"encoding/json"
	"golib/system_solution/chat/model"
	"strconv"
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

func (c *chatService) formatUserOnlineStatusKey(userId int64) string {
	return "chat:online:user" + strconv.FormatInt(userId, 10)
}

type onlineUserInfo struct {
	ServiceHost string `json:"server_host"`
}

func newOnlineUserInfo(serviceHost string) *onlineUserInfo {
	return &onlineUserInfo{ServiceHost: serviceHost}
}

func (oui *onlineUserInfo) String() string {
	outBytes, err := json.Marshal(oui)
	if err != nil {
		return ""
	}
	return string(outBytes)
}

// PushMessage
//
//	@Description:
//		1.检查用户是否本机在线，是的话则直接推送，不是的话则通过ETCD获取用户所在服务器，然后推送
//		2.如果用户不在线则，则看是否有屏蔽消息，没有屏蔽消息则不推入消息中心，等待用户拉取
func (cs *chatService) PushMessage(ctx context.Context, message *model.Message) error {
	channel := ChannelManager().GetChannel(message.ToUser)
	if channel != nil {
		err := channel.PushMessage(message)
		if err != nil {
			return err
		}
	}

	return nil
}

// UserOnline
//
//	@Description: 用户上线后将信息维护到ETCD中，通过ETCD可以获取用户在哪个服务器上
func (cs *chatService) UserOnline(ctx context.Context, userId int64) error {
	return nil
}

// UserOffline
//
//	@Description: TODO 用户下线后将信息从ETCD中删除
func (cs *chatService) UserOffline(ctx context.Context, userId int64) error {
	return nil
}
