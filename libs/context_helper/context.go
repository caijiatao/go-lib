package context_helper

import (
	"context"
	"reflect"
	"unsafe"
)

type iface struct {
	itab, data uintptr
}

type valueCtx struct {
	context.Context
	key, val interface{}
}

func GetKeyValues(ctx context.Context) map[interface{}]interface{} {
	m := make(map[interface{}]interface{})
	getKeyValue(ctx, m)
	return m
}

func getKeyValue(ctx context.Context, m map[interface{}]interface{}) {

	rtType := reflect.TypeOf(ctx).String()

	// 遍历到顶级类型，直接过滤
	if rtType == "*context.emptyCtx" {
		return
	}

	ictx := *(*iface)(unsafe.Pointer(&ctx))
	if ictx.data == 0 {
		return
	}
	valCtx := (*valueCtx)(unsafe.Pointer(ictx.data))
	if valCtx != nil && valCtx.key != nil && valCtx.val != nil {
		m[valCtx.key] = valCtx.val
	}
	getKeyValue(valCtx.Context, m)
}

func DeepCopy(ctx context.Context) context.Context {
	keyValues := GetKeyValues(ctx)
	newCtx := context.Background()
	for k, v := range keyValues {
		newCtx = context.WithValue(newCtx, k, v)
	}
	return newCtx
}
