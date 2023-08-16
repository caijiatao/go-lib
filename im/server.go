package im

import (
	"github.com/gorilla/websocket"
	"net/http"
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

func Register(w http.ResponseWriter, r *http.Request, key string) error {
	client, err := connect(w, r)
	if err != nil {
		return err
	}
	client.Key = key
	err = globalClientManager.Register(client)
	if err != nil {
		return err
	}
	return nil
}

func connect(w http.ResponseWriter, r *http.Request) (*Client, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	client := NewClient(conn)

	go client.write()
	go client.read()

	return client, nil
}
