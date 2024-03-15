package chat

import (
	"context"
	"github.com/gorilla/websocket"
)

type Channel struct {
	userId int64
	conn   *websocket.Conn
	send   chan []byte
}

func (c *Channel) sendLoop(ctx context.Context) {

}

func (c *Channel) recvLoop(ctx context.Context) {

}

type ChannelManager struct {
	userId2Channel map[int64]Channel
}

func (m *ChannelManager) AddChannel(userId int64, channel Channel) {

}
