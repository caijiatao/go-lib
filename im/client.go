package im

import (
	"github.com/gorilla/websocket"
	"golib/logger"
	"time"
)

const (
	writeWait = time.Second * 10
)

type Client struct {
	conn *websocket.Conn

	send chan []byte
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		conn: conn,
	}
}

func (c *Client) write() {
	for {
		select {
		case message, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				err := c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					logger.Errorf("write close message error: %v", err)
				}
				return
			}

			writer, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			n, err := writer.Write(message)
			if err != nil {
				logger.Errorf("write message error: %v", err)
				continue
			}
			logger.Infof("write message success, len: %d", n)
			err = writer.Close()
			if err != nil {
				logger.Errorf("close writer error: %v", err)
				return
			}
		}
	}
}

func (c *Client) read() {
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			logger.Errorf("read message error: %v", err)
			return
		}
		logger.Infof("read message success, message: %s", string(message))
	}
}
