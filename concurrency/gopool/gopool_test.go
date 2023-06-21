package gopool

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	initManager()
	initTask()
	initWorker()
	m.Run()
}

func execFunc() {
	// 模拟耗时操作
	time.Sleep(10 * time.Millisecond)
	fmt.Println(WorkCount())
}

func TestGoPool(t *testing.T) {
	execTimes := 1000000
	var wg sync.WaitGroup
	wg.Add(execTimes)
	for i := 0; i < execTimes; i++ {
		Go(func() {
			execFunc()
			wg.Done()
		})
	}
	wg.Wait()
}
