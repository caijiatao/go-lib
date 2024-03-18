package logger

import (
	"context"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golib/libs/osutil"
	"os"
	"path"
	"strings"
	"time"
)

type zapLogger struct {
	l     *zap.SugaredLogger
	level logLevel
}

func (z *zapLogger) Debug(v ...interface{}) {
	if z.level <= logDebug {
		z.l.Debug(v...)
	}
}

func (z *zapLogger) Debugf(format string, v ...interface{}) {
	if z.level <= logDebug {
		z.l.Debugf(format, v...)
	}
}

func (z *zapLogger) CtxDebugf(ctx context.Context, format string, v ...interface{}) {
	if z.level <= logDebug {
		z.l.With(zap.String(string(TraceIDKey), CtxTraceID(ctx))).Debugf(format, v...)
	}
}

func (z *zapLogger) Info(v ...interface{}) {
	if z.level <= logInfo {
		z.l.Info(v...)
	}
}

func (z *zapLogger) Infof(format string, v ...interface{}) {
	if z.level <= logInfo {
		z.l.Infof(format, v...)
	}
}

func (z *zapLogger) CtxInfof(ctx context.Context, format string, v ...interface{}) {
	if z.level <= logInfo {
		z.l.With(zap.String(string(TraceIDKey), CtxTraceID(ctx))).Infof(format, v...)
	}
}

func (z *zapLogger) Warn(v ...interface{}) {
	if z.level <= logWarn {
		z.l.Warn(v...)
	}
}

func (z *zapLogger) Warnf(format string, v ...interface{}) {
	if z.level <= logWarn {
		z.l.Warnf(format, v...)
	}
}

func (z *zapLogger) CtxWarnf(ctx context.Context, format string, v ...interface{}) {
	if z.level <= logWarn {
		z.l.With(zap.String(string(TraceIDKey), CtxTraceID(ctx))).Warnf(format, v...)
	}
}

func (z *zapLogger) Error(v ...interface{}) {
	if z.level <= logErr {
		z.l.Error(v...)
	}
}

func (z *zapLogger) Errorf(format string, v ...interface{}) {
	if z.level <= logErr {
		z.l.Errorf(format, v...)
	}
}

func (z *zapLogger) CtxErrorf(ctx context.Context, format string, v ...interface{}) {
	if z.level <= logErr {
		z.l.With(zap.String(string(TraceIDKey), CtxTraceID(ctx))).Errorf(format, v...)
	}
}

func newZapLogger() *zapLogger {

	maxAge := time.Duration(config.MaxAge) * 24 * time.Hour
	rotationTime := time.Duration(config.RotationTime) * time.Hour

	// 创建日志存放目录
	if err := osutil.CreateDirIfNotExist(config.Path); err != nil {
		panic(err)
	}
	logPath := path.Join(config.Path, config.ProjectName)

	// error日志文件配置
	errWriter, err := rotatelogs.New(
		logPath+"_err_%Y-%m-%d.log",
		rotatelogs.WithMaxAge(maxAge),
		rotatelogs.WithRotationTime(rotationTime),
	)
	if err != nil {
		panic(err)
	}

	// info日志文件配置
	infoWriter, err := rotatelogs.New(
		logPath+"_info_%Y-%m-%d.log",
		rotatelogs.WithMaxAge(maxAge),
		rotatelogs.WithRotationTime(rotationTime),
	)
	if err != nil {
		panic(err)
	}

	// 优先级设置
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl > zapcore.WarnLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.DebugLevel
	})

	// 控制台输出设置
	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleEncoderConfig := zap.NewDevelopmentEncoderConfig()
	consoleEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoderConfig.EncodeTime = timeEncoder
	consoleEncoderConfig.EncodeCaller = customCallerEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderConfig)

	// 文件输出设置
	errorCore := zapcore.AddSync(errWriter)
	infoCore := zapcore.AddSync(infoWriter)
	fileEncodeConfig := zap.NewProductionEncoderConfig()
	fileEncodeConfig.EncodeTime = timeEncoder
	fileEncodeConfig.EncodeCaller = customCallerEncoder
	fileEncoder := zapcore.NewConsoleEncoder(fileEncodeConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, errorCore, highPriority),
		zapcore.NewCore(fileEncoder, infoCore, lowPriority),
		zapcore.NewCore(consoleEncoder, consoleDebugging, zapcore.DebugLevel),
	)

	// 显示行号
	caller := zap.AddCaller()

	development := zap.Development()
	zapLog := zap.New(core, caller, development)

	// 替换全局日志
	zap.ReplaceGlobals(zapLog)

	// 将系统输出重定向到zap中，保证所有出现异常均能打印到文件中
	if _, err := zap.RedirectStdLogAt(zapLog, zapcore.ErrorLevel); err != nil {
		panic(err)
	}

	return &zapLogger{
		l:     zap.L().WithOptions(zap.AddCallerSkip(2)).Sugar(),
		level: levelMap()[config.Level],
	}

}

// customCallerEncoder 自定义打印路径，减少输出日志打印路径长度，根据输入项目名进行减少
func customCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	str := caller.String()
	index := strings.Index(str, config.ProjectName)
	if index == -1 {
		enc.AppendString(caller.FullPath())
	} else {
		index = index + len(config.ProjectName) + 1
		enc.AppendString(str[index:])
	}
}

// timeEncoder 格式化日志时间，官方的不好看
func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}
