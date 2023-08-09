package logger

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
)

// GetGinRequestBody 获取请求 body
func GetGinRequestBody(c *gin.Context) []byte {
	// 获取请求 body
	var requestBody []byte
	if c.Request.Body != nil {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.Error(err)
		} else {
			requestBody = body
			// body 被 read 、 bind 之后会被置空，需要重置
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		}
	}
	return requestBody
}
