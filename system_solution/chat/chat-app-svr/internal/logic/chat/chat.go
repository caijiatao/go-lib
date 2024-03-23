package service

import (
	"chat-app-svr/internal/config"
	"chat-app-svr/internal/model"
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"golib/libs/etcd_helper"
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

func (cs *chatService) formatUserOnlineStatusKey(userId int64) string {
	return "/chat/online/user" + strconv.FormatInt(userId, 10)
}

type onlineUserInfo struct {
	ServerID string `json:"server_id"`
}

func newOnlineUserInfo(serverID string) *onlineUserInfo {
	return &onlineUserInfo{
		ServerID: serverID,
	}
}

func (oui *onlineUserInfo) String() string {
	outBytes, err := json.Marshal(oui)
	if err != nil {
		return ""
	}
	return string(outBytes)
}

func (cs *chatService) checkMessage(message *model.Message) error {
	if message == nil {
		return errors.New("message is empty")
	}
	if message.ToUser == 0 {
		return errors.New("to_user is empty")
	}
	if message.FromUser == 0 {
		return errors.New("from_user is empty")
	}
	if len(message.Content) == 0 {
		return errors.New("content is empty")
	}
	return nil
}

// PushMessage
//
//	@Description:
//		1.检查用户是否本机在线，是的话则直接推送，不是的话则通过ETCD获取用户所在服务器，然后推送
//		2.如果用户不在线则，则看是否有屏蔽消息，没有屏蔽消息则不推入消息中心，等待用户拉取
func (cs *chatService) PushMessage(ctx context.Context, message *model.Message) error {
	err := cs.checkMessage(message)
	if err != nil {
		return err
	}
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
	userOnlineInfo := newOnlineUserInfo(config.Conf().GetServerID())
	putResp, err := etcd_helper.Put(ctx, cs.formatUserOnlineStatusKey(userId), userOnlineInfo.String())
	if err != nil {
		return err
	}
	if putResp.Header.Revision == 0 {
		return errors.New("user online failed")
	}
	return nil
}

// UserOffline
//
//	@Description: 用户下线后将信息从ETCD中删除
func (cs *chatService) UserOffline(ctx context.Context, userId int64) error {
	_, err := etcd_helper.Delete(ctx, cs.formatUserOnlineStatusKey(userId))
	if err != nil {
		return err
	}
	return nil
}
