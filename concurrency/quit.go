package concurrency

import (
	"golib/concurrency/gopool"
	"golib/logger"
	"golib/util"
	"os"
	"os/signal"
	"syscall"
)

func GracefulQuit(cancel func()) {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	callerName := util.GetCallerFunctionName()
	gopool.Go(func() {
		osCall := <-quit
		logger.Infof("system call: %+v, callerName:%s", osCall, callerName)
		cancel()
	})
}
