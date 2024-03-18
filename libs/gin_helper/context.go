package gin_helper

import (
	"github.com/gin-gonic/gin"
	"golib/libs/etcd_helper"
	"golib/libs/orm"
)

func Context() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 在中间件中使用 context.WithValue 设置值
		ctx := orm.BindContext(c.Request.Context())
		ctx = etcd_helper.BindContext(ctx)

		c.Request = c.Request.WithContext(ctx)

		// 继续处理下一个中间件或处理函数
		c.Next()
	}
}
