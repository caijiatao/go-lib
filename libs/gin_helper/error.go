package gin_helper

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

type MyError struct {
	Code       int    `json:"code"`
	Msg        string `json:"msg"`
	ServerTime uint
	Data       interface{}
}

var (
	PARAM_ERROR    = NewError(-1, "参数错误")
	NAME_DUPLICATE = NewError(-1, "名称重复")

	ERROR           = NewError(-1, "操作失败")
	INNER_ERROR     = NewError(-1, "服务异常")
	LIMIT_ERROR     = NewError(-1, "服务器限流，请稍后重试")
	PERMISSION_DENY = NewError(-1, "权限拒绝")

	LOGIN_UNKNOWN    = NewError(202, "该账号未注册")
	LOGIN_ERROR      = NewError(203, "登录失败，密码输入错误，剩余可登录次数%s次，连续输错5次将被锁定账号30分钟")
	LOGIN_LOCK       = NewError(204, "登录失败，登录次数达到5次，账号锁定30分钟，还有 %d 分钟 %d 秒可再次登录")
	TOKEN_TIMEOUT    = NewError(401, "token超时失效")
	UNAUTHORIZED     = NewError(402, "认证失败，请登录")
	SMSCODE_ERROR    = NewError(403, "验证码错误，请输入正确的验证码")
	SMSCODE_LIMIT    = NewError(404, "今日已达验证码最大发送次数，请明日重试")
	SAME_PASSWORD    = NewError(405, "新密码与原始密码重复，请不要设置相同密码")
	ERR_PASSWORD     = NewError(406, "密码错误，请输入正确的密码")
	ILLEGAL_PASSWORD = NewError(407, "密码非法，请按要求设置正确密码")

	MAX_NAMESPACE = NewError(501, "当前已达到体验Demo人数上限，若需继续体验，请提交申请试用")

	CONNECT_FAILED = NewError(-1, "数据库连接失败")

	MAPPING_MISSING        = NewError(-1, "映射字段缺失")
	COLUMN_TYPE_UNMATCH    = NewError(-1, "不同类型字段必须配置转化脚本")
	SOURCE_TABLE_DUPLICATE = NewError(-1, "源表重复添加")
	SOURCE_IN_USE          = NewError(-1, "该数据源已被添加或关联，不可删除")
)

func (e *MyError) Error() string {
	return e.Msg
}

func (e *MyError) String() string {
	errString, err := json.Marshal(e)
	if err != nil {
		return e.Error()
	}
	return string(errString)
}

func GetErrString(err error) string {
	if err == nil {
		return ""
	}
	if myErr, ok := err.(*MyError); ok {
		return myErr.String()
	}
	return err.Error()
}

func GetMyErr(err error) *MyError {
	if err == nil {
		return nil
	}
	err = errors.Cause(err)
	if myErr, ok := err.(*MyError); ok {
		return myErr
	}
	return nil
}

func NewError(code int, msg string) *MyError {
	return &MyError{
		Msg:        msg,
		Code:       code,
		ServerTime: uint(time.Now().Unix()),
	}
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // 先调用c.Next()执行后面的中间件
		// 所有中间件及router处理完毕后从这里开始执行
		// 检查c.Errors中是否有错误
		for _, e := range c.Errors {
			err := e.Err
			// 若是自定义的错误则将code、msg返回
			if myErr, ok := err.(*MyError); ok {
				c.JSON(http.StatusOK, gin.H{
					"code":       myErr.Code,
					"msg":        myErr.Msg,
					"serverTime": myErr.ServerTime,
					"data":       myErr.Data,
				})
			} else {
				// 若非自定义错误则返回详细错误信息err.Error()
				// 比如save session出错时设置的err
				c.JSON(http.StatusOK, gin.H{
					"code":       -1,
					"msg":        err.Error(),
					"serverTime": time.Now().Unix(),
					"data":       nil,
				})
			}
			return // 检查一个错误就行
		}
	}
}
