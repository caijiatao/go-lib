package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	err := App().Init(gin.Default()).Run()
	if err != nil {
		panic(err)
	}
}
