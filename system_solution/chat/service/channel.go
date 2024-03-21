package service

import (
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"golib/libs/etcd_helper"
	"golib/libs/logger"
	"golib/libs/orm"
	"golib/system_solution/chat/model"
	"sync"
)

type Channel struct {
	userId int64
	conn   *websocket.Conn
	send   chan []byte
}

func NewChannel(userId int64, conn *websocket.Conn, send chan []byte) *Channel {
	return &Channel{userId: userId, conn: conn, send: send}
}

func (c *Channel) PushMessage(message *model.Message) error {
	ms, err := json.Marshal(message)
	if err != nil {
		return err
	}
	c.send <- ms
	return nil
}

func (c *Channel) SendLoop() {
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				return
			}
			err := c.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				return
			}
		}
	}
}

func (c *Channel) RecvLoop() {
	ctx := orm.BindContext(context.Background())
	for {
		messageType, message, err := c.conn.ReadMessage()
		if err != nil {
			return
		}

		switch messageType {
		case websocket.PingMessage:
			err = c.conn.WriteMessage(websocket.PongMessage, nil)
		case websocket.TextMessage:
			m := &model.Message{}
			err = json.Unmarshal(message, m)
			if err != nil {
				return
			}
			m.FromUser = c.userId
			err = ChatService().PushMessage(ctx, m)
			if err != nil {
				logger.CtxErrorf(ctx, "push message error: %v", err)
			}
		}
		if err != nil {
			_ = c.Close()
			return
		}
	}
}

func (c *Channel) Close() error {
	ChannelManager().RemoveChannel(c)
	close(c.send)
	err := c.conn.Close()
	if err != nil {
		logger.CtxErrorf(nil, "close connection error: %v", err)
		return err
	}
	return nil
}

var (
	channelManager         *manager
	channelManagerInitOnce = sync.Once{}
)

func ChannelManager() *manager {
	channelManagerInitOnce.Do(func() {
		channelManager = &manager{
			userId2Channel: make(map[int64]*Channel),
		}
		channelManager.run()
	})
	return channelManager
}

type manager struct {
	sync.RWMutex
	userId2Channel map[int64]*Channel

	unregister chan *Channel
	register   chan *Channel
}

func (m *manager) run() {
	go m.registerLoop()
	go m.unregisterLoop()
}

func (m *manager) registerLoop() {
	ctx := context.Background()
	ctx = orm.BindContext(ctx)
	ctx = etcd_helper.BindContext(ctx)
	for {
		select {
		case c, ok := <-m.register:
			if !ok {
				return
			}
			m.userId2Channel[c.userId] = c
			err := ChatService().UserOnline(ctx, c.userId)
			if err != nil {
				logger.CtxErrorf(ctx, "user online error: %v", err)
			}
		}
	}
}

func (m *manager) unregisterLoop() {
	ctx := orm.BindContext(context.Background())
	for {
		select {
		case c, ok := <-m.unregister:
			if !ok {
				return
			}
			delete(m.userId2Channel, c.userId)
			err := ChatService().UserOffline(ctx, c.userId)
			if err != nil {
				logger.CtxErrorf(nil, "user offline error: %v", err)
			}
		}
	}
}

func (m *manager) AddChannel(channel *Channel) {
	m.register <- channel
}

func (m *manager) RemoveChannel(channel *Channel) {
	m.unregister <- channel
}

func (m *manager) GetChannel(userId int64) *Channel {
	m.RLock()
	defer m.RUnlock()
	return m.userId2Channel[userId]
}

func (m *manager) PushMessage(message *model.Message) {
	c := m.GetChannel(message.ToUser)
	if c != nil {
		err := c.PushMessage(message)
		if err != nil {
			logger.CtxErrorf(nil, "push message error: %v", err)
		}
		return
	}

	// 用户没在本机，则转发到其他机器

	// 用户没在线，则存储到数据库
}
