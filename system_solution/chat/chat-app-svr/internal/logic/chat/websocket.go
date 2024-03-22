package service

import (
	"github.com/gorilla/websocket"
	"golib/system_solution/chat/service"
	"log"
	"net/http"
	"strconv"
)

func ServeWs(w http.ResponseWriter, r *http.Request) {
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// 在这里修改为允许特定来源
			// 比如允许所有来源：
			// return true
			// 或者允许特定的来源：
			// return r.Header.Get("Origin") == "http://your.allowed.origin"
			return true
		},
	}).Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	if err != nil {
		return
	}
	// parse user id
	query := r.URL.Query()
	userID := query.Get("user_id")

	if userID == "" {
		return
	}
	// user id to int64
	userIDInt, err := strconv.ParseInt(userID, 10, 64)
	channel := service.NewChannel(userIDInt, conn, make(chan []byte, 256))

	go channel.SendLoop()
	go channel.RecvLoop()

	service.ChannelManager().AddChannel(channel)
}
