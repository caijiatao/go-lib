package logger

import (
	"testing"
)

func TestTraceId(t *testing.T) {
	ctx := CtxWithTraceId(nil, "")
	defaultLogger.L.CtxErrorf(ctx, "test trace id:%s", "tttt")
	defaultLogger.L.CtxDebugf(ctx, "test trace id:%s", "tttt")
	defaultLogger.L.CtxInfof(ctx, "test trace id:%s", "tttt")
	defaultLogger.L.CtxWarnf(ctx, "test trace id:%s", "tttt")
}
