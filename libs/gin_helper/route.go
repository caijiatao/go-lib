package gin_helper

import (
	"github.com/gin-gonic/gin"
	"reflect"
)

type IRoute interface {
	RegisterRoutes(r gin.IRouter)
}

func RegisterAllRoutes(obj interface{}, r gin.IRouter) {
	objValue := reflect.ValueOf(obj)

	if objValue.Kind() != reflect.Pointer {
		return
	}

	for i := 0; i < objValue.Elem().NumField(); i++ {
		field := objValue.Elem().Field(i)
		if field.Kind() == reflect.Pointer {
			route, ok := field.Interface().(IRoute)
			if ok {
				route.RegisterRoutes(r)
			}
		}
	}
}
