package im

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Server struct {
}

func InitWebSocket() (err error) {
	return nil
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func Connect(c *gin.Context, key string) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	client := NewClient(conn)

	go client.write()
	go client.read()
}
