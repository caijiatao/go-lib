package concurrency

import (
	"fmt"
	"sync"
	"time"
)

type BoundedFrequencyRunner struct {
	sync.Mutex

	// 主动触发
	run chan struct{}

	// 定时器限制
	timer *time.Timer

	// 真正执行的逻辑
	fn func()
}

func NewBoundedFrequencyRunner(fn func()) *BoundedFrequencyRunner {
	return &BoundedFrequencyRunner{
		run:   make(chan struct{}, 1),
		fn:    fn,
		timer: time.NewTimer(0),
	}
}

// Run 触发执行 ,这里只能够写入一个信号量，多余的直接丢弃，不会阻塞，这里也可以根据自己的需要增加排队的个数
func (b *BoundedFrequencyRunner) Run() {
	select {
	case b.run <- struct{}{}:
		fmt.Println("写入信号量成功")
	default:
		fmt.Println("已经触发过一次，直接丢弃信号量")
	}
}

func (b *BoundedFrequencyRunner) Loop() {
	b.timer.Reset(time.Second * 1)
	for {
		select {
		case <-b.run:
			fmt.Println("run 信号触发")
			b.tryRun()
		case <-b.timer.C:
			fmt.Println("timer 触发执行")
			b.tryRun()
		}
	}
}

func (b *BoundedFrequencyRunner) tryRun() {
	b.Lock()
	defer b.Unlock()
	// 可以增加限流器等限制逻辑
	b.timer.Reset(time.Second * 1)
	b.fn()
}
