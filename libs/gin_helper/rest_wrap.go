package gin_helper

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type RestHandler func(ctx *gin.Context) (interface{}, error)

type Response struct {
	Code       int         `json:"code"`
	Message    string      `json:"msg"`
	ServerTime uint        `json:"ServerTime"`
	Data       interface{} `json:"data,omitempty"`
}

// SendSuccessResp 返回成功请求
func SendSuccessResp(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:       0,
		Message:    "操作成功",
		ServerTime: uint(time.Now().Unix()),
		Data:       data,
	})
}

//func RestWrapper(h RestHandler) func(ctx *gin.Context) {
//	return func(ctx *gin.Context) {
//		// 从panic中恢复
//		defer func() {
//			e := recover()
//			if e != nil {
//				err := goasync.PanicErrHandler(e)
//				ctx.JSON(http.StatusOK, err.Error())
//				return
//			}
//		}()
//		data, err := h(ctx)
//		if err != nil {
//			//err = errors.Cause(err)
//			//ctx.JSON(http.StatusOK, Response{
//			//	Code:       0,
//			//	Message:    err.Error(),
//			//	ServerTime: uint(time.Now().Unix()),
//			//})
//			// 若是自定义的错误则将code、msg返回
//			if myErr, ok := err.(*errors.MyError); ok {
//				ctx.JSON(http.StatusOK, gin.H{
//					"code":       myErr.Code,
//					"msg":        myErr.Msg,
//					"serverTime": myErr.ServerTime,
//					"data":       myErr.Data,
//				})
//			} else {
//				// 若非自定义错误则返回详细错误信息err.Error()
//				// 比如save session出错时设置的err
//				ctx.JSON(http.StatusOK, gin.H{
//					"code":       -1,
//					"msg":        "服务异常",
//					"serverTime": time.Now().Unix(),
//					"data":       nil,
//				})
//			}
//		} else {
//			ctx.JSON(http.StatusOK, Response{
//				Code:       0,
//				Message:    "success",
//				ServerTime: uint(time.Now().Unix()),
//				Data:       data,
//			})
//		}
//	}
//}
