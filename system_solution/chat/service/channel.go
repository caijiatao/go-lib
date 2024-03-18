package service

import (
	"context"
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

func (c *Channel) SendLoop(ctx context.Context) {
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
		case <-ctx.Done():
			return
		}
	}
}

func (c *Channel) RecvLoop(ctx context.Context) {
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			return
		}
		m := Message{}
		err = json.Unmarshal(message, &m)
		if err != nil {
			return
		}
	}
}

type ChannelManager struct {
	sync.Mutex
	userId2Channel map[int64]Channel
}

func (m *ChannelManager) AddChannel(userId int64, channel Channel) {
	m.Lock()
	defer m.Unlock()
	m.userId2Channel[userId] = channel
}
