package chat

import (
	"chat-app-svr/internal/svc"
	"chat-app-svr/rpc/user/user"
	"context"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type WebsocketHandler struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWebsocketHandler(ctx context.Context, svcCtx *svc.ServiceContext) *WebsocketHandler {
	_ = NewChannelManager(svcCtx)
	return &WebsocketHandler{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (wsHandler *WebsocketHandler) GetUserIdByRequest(r *http.Request) (int64, error) {
	token := r.Header.Get("token")
	authReply, err := wsHandler.svcCtx.User.Auth(wsHandler.ctx, &user.AuthRequest{
		Token: token,
	})
	if err != nil {
		return 0, err
	}
	return authReply.UserInfo.UserId, nil
}

func (wsHandler *WebsocketHandler) ServeWs(w http.ResponseWriter, r *http.Request) {
	userID, err := wsHandler.GetUserIdByRequest(r)
	if err != nil {
		return
	}

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

	channel := NewChannel(userID, conn, make(chan []byte, 256), wsHandler.svcCtx)

	go channel.SendLoop()
	go channel.RecvLoop()

	ChannelManager().AddChannel(channel)
}
