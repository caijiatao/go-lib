package concurrency

import (
	"context"
	"fmt"
	"golib/libs/logger"
	"golib/libs/util"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func GracefulQuit(cancel func()) {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT)
	callerName := util.GetCallerFunctionName()
	go func() {
		osCall := <-quit
		logger.Infof("system call: %+v, callerName:%s", osCall, callerName)
		cancel()
	}()
}

type GracefulShutdown interface {
	Shutdown(ctx context.Context) error
}

// ShutdownFunc 能够将函数方法转化成符合类型的接口
type ShutdownFunc func(ctx context.Context) error

func (f ShutdownFunc) Shutdown(ctx context.Context) error {
	return f(ctx)
}

var (
	shutdownerInitOnce sync.Once
	globalShutdowner   *ApplicationShutdowner
)

type ApplicationShutdowner struct {
	waitTimeout time.Duration

	shutdowns []GracefulShutdown
}

func NewApplicationShutdowner() *ApplicationShutdowner {
	// 只有在第一个注册的时候会实例化
	shutdownerInitOnce.Do(func() {
		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT)
		callerName := util.GetCallerFunctionName()
		go func() {
			osCall := <-quit
			fmt.Println("start shutdown")
			logger.Infof("system call: %+v, callerName:%s", osCall, callerName)
			globalShutdowner.Shutdown()
			// 执行完了，退出程序
			os.Exit(0)
		}()

		globalShutdowner = &ApplicationShutdowner{
			waitTimeout: time.Second * 15, // 默认15秒，不暴露修改
		}
	})
	return globalShutdowner
}

func RegisterShutdown(shutdown GracefulShutdown) {
	shutdowner := NewApplicationShutdowner()
	shutdowner.shutdowns = append(shutdowner.shutdowns, shutdown)
}

func (ss *ApplicationShutdowner) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), ss.waitTimeout)
	defer cancel()
	for _, shutdown := range ss.shutdowns {
		err := shutdown.Shutdown(ctx)
		if err != nil {
			return
		}
	}
}
