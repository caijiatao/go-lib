package logger

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"golib/libs/goasync"
	"io"
	"net/http"
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

func MiddlewareLoggerInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		ctx = CtxWithTraceId(ctx, "")
		c.Request = c.Request.WithContext(ctx)
		blw := &CustomResponseWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		defer func() {
			r := recover()
			_ = goasync.PanicErrHandler(r)
		}()
		CtxInfof(c.Request.Context(), "WEB_BEFORE---------------------------------------------------")
		CtxInfof(c.Request.Context(), "IP: %s", c.Request.Host)
		CtxInfof(c.Request.Context(), "URL: %s", c.Request.URL)
		CtxInfof(c.Request.Context(), "HTTP_METHOD: %s", c.Request.Method)
		// PUT或者POST方法时打印body
		if (c.Request.Method == http.MethodPut || c.Request.Method == http.MethodPost) &&
			c.Request.Header.Get("Content-Type") == "application/json" {
			data, err := io.ReadAll(c.Request.Body)
			if err != nil {
				CtxInfof(c.Request.Context(), "read request body failed,err = %s.", err)
				return
			}
			CtxInfof(c.Request.Context(), "ARGS: %s", string(data))
			// 重新赋值
			c.Request.Body = io.NopCloser(bytes.NewBuffer(data))
		}
		c.Next()
		if c.Request.Method != http.MethodGet {
			CtxInfof(c.Request.Context(), "Resp: %s", blw.body.String())
		}
		CtxInfof(c.Request.Context(), "WEB_AFTER----------------------------------------------------")
	}
}

type CustomResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w CustomResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w CustomResponseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}
