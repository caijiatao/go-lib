package main

import (
	"fmt"
	//"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	//ctx := gin.Context{}
	db := gorm.DB{}
	fmt.Println(db)
	//fmt.Println(ctx)
}
