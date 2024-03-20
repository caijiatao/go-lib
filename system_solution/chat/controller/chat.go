package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"golib/libs/gin_helper"
	"golib/system_solution/chat/service"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type Controller struct{}

var (
	chatController         *Controller
	chatControllerInitOnce = sync.Once{}
)

func ChatController() *Controller {
	chatControllerInitOnce.Do(func() {
		chatController = &Controller{}
	})
	return chatController
}

func (c *Controller) RegisterRoutes(r gin.IRouter) {
	r.Use(gin_helper.CorsMiddleware())

	r.GET("/", c.HomePage)
	r.GET("/ws/chat", c.Chat)
}

func (c *Controller) HomePage(ctx *gin.Context) {
	log.Println(ctx.Request.URL)
	if ctx.Request.URL.Path != "/" {
		http.Error(ctx.Writer, "Not found", http.StatusNotFound)
		return
	}
	http.ServeFile(ctx.Writer, ctx.Request, "C:\\project\\github\\go-lib\\system_solution\\chat\\home.html")
}

func (c *Controller) Chat(ctx *gin.Context) {
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// 在这里修改为允许特定来源
			// 比如允许所有来源：
			// return true
			// 或者允许特定的来源：
			// return r.Header.Get("Origin") == "http://your.allowed.origin"
			return true
		},
	}).Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return
	}
	// parse user id
	userID := ctx.Query("user_id")
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
