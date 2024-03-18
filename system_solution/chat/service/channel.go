package service

import (
	"encoding/json"
	"github.com/gorilla/websocket"
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
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			return
		}
		m := &Message{}
		err = json.Unmarshal(message, m)
		if err != nil {
			return
		}
		m.FromUser = c.userId
		m.ToUser = c.userId // 重新发给自己

		ChannelManager().PushMessage(m)
	}
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
	})
	return channelManager
}

type manager struct {
	sync.Mutex
	userId2Channel map[int64]*Channel
}

func (m *manager) AddChannel(channel *Channel) {
	m.Lock()
	defer m.Unlock()
	m.userId2Channel[channel.userId] = channel
}

func (m *manager) RemoveChannel(userId int64) {
	m.Lock()
	defer m.Unlock()
	delete(m.userId2Channel, userId)
}

func (m *manager) PushMessage(message *Message) {
	m.Lock()
	channel, ok := m.userId2Channel[message.ToUser]
	m.Unlock()

	if ok {
		ms, err := json.Marshal(message)
		if err != nil {
			return
		}
		channel.send <- ms
	}
	// 用户没在本机，则转发到其他机器

	// 用户没在线，则存储到数据库
}
