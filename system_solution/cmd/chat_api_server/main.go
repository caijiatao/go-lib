package main

import (
	"github.com/gin-gonic/gin"
	"golib/system_solution/chat_api_server"
	"golib/system_solution/cmd/chat_api_server/config"
)

func main() {
	chat_api_server.NewServer(config.Conf())

	err := App().Init(gin.Default()).Run()
	if err != nil {
		panic(err)
	}
}
