package chat

import (
	"context"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

type Server struct{}

var (
	server         *Server
	serverInitOnce = sync.Once{}
)

func GetServer() *Server {
	serverInitOnce.Do(func() {
		server = &Server{}
	})
	return server
}

func (c *Server) RegisterRoutes() {
	http.HandleFunc("/", c.HomePage)
	http.HandleFunc("/ws/chat", c.Chat)
}

func (c *Server) HomePage(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "C:\\project\\github\\go-lib\\system_solution\\chat\\home.html")
}

func (c *Server) Chat(writer http.ResponseWriter, request *http.Request) {
	conn, err := (&websocket.Upgrader{}).Upgrade(writer, request, nil)
	if err != nil {
		return
	}

	ctx := context.Background()
	channel := Channel{
		userId: 0,
		conn:   conn,
		send:   make(chan []byte, 100),
	}

	go channel.sendLoop(ctx)

	go channel.recvLoop(ctx)
}
