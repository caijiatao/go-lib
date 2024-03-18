package gin_helper

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") //请求头部
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "POST,GET,OPTIONS,PUT,DELETE,UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin,Content-Type,Accept,Authorization,Content-Length,Token")
			c.Header("Access-Control-Expose-Headers", "Access-Control-Allow-Headers,Token")
			c.Header("Access-Control-Max-Age", "86400") // 可选
		}

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}

		c.Next()
	}
}
