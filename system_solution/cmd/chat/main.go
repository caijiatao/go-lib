package main

import (
	"github.com/gin-gonic/gin"
	"golib/system_solution/chat"
	"golib/system_solution/cmd/chat/config"
)

func main() {
	chat.NewServer(config.Conf())

	err := App().Init(gin.Default()).Run()
	if err != nil {
		panic(err)
	}
}
