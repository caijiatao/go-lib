package main

import (
	"github.com/gin-gonic/gin"
)

func watchChange(ctx *gin.Context) {

}

func main() {
	engine := gin.Default()

	engine.POST("/watch_change", watchChange)

	err := engine.Run(":8080")
	if err != nil {
		panic(err)
	}
}
