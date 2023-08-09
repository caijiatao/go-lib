package logger

import (
	"context"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/xid"
)

// Ctxkey context key 类型
type Ctxkey string

var (
	TraceIDKey    Ctxkey = "trace_id"
	TraceIDPrefix        = "logging_"
)

func WithTraceId(ctx context.Context) context.Context {
	return CtxWithTraceId(ctx, "")
}

func CtxWithTraceId(ctx context.Context, traceId string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if traceId == "" {
		traceId = CtxTraceID(ctx)
	}
	ctx = context.WithValue(ctx, TraceIDKey, traceId)
	return ctx
}

func CtxTraceID(c context.Context) string {
	if c == nil {
		c = context.Background()
	}
	// 从gin的请求去获取trace id
	if gc, ok := c.(*gin.Context); ok {
		if traceID := gc.GetString(string(TraceIDKey)); traceID != "" {
			return traceID
		}
		if traceID := gc.Query(string(TraceIDKey)); traceID != "" {
			return traceID
		}
		if traceID := jsoniter.Get(GetGinRequestBody(gc), string(TraceIDKey)).ToString(); traceID != "" {
			return traceID
		}
	}
	// get from go context
	traceID := c.Value(TraceIDKey)
	if traceID != nil {
		return traceID.(string)
	}
	return TraceIDPrefix + xid.New().String()
}
