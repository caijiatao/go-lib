package concurrency

import (
	"sync"
	"testing"
	"time"
)

func Test_write(t *testing.T) {
	cond := sync.NewCond(&sync.Mutex{})
	// 启动协程去读取，还未写入则不能读取成功
	go read(cond)

	go read(cond)

	time.Sleep(time.Second * 1)

	// 写入数据，完成后就能够读取成功
	write(cond)

	time.Sleep(time.Second * 3)
}
