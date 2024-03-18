package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"golib/system_solution/chat/service"
	"golib/system_solution/user/middleware"
	"log"
	"net/http"
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
	conn, err := (&websocket.Upgrader{}).Upgrade(ctx.Writer, ctx.Request, nil)
	user := middleware.AuthMiddleware(ctx).GetUser()
	if err != nil {
		return
	}
	channel := service.NewChannel(user.UserID, conn, make(chan []byte, 256))

	go channel.SendLoop(ctx)

	go channel.RecvLoop(ctx)
}
