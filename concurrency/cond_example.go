package concurrency

import (
	"fmt"
	"sync"
	"time"
)

var done bool

func read(c *sync.Cond) {
	c.L.Lock()
	// 写入未完成之前先陷入等待
	for !done {
		fmt.Println("read func wait")
		// 内部有unlock，等待被通知后重新唤起，所以临界资源仍然是安全的
		c.Wait()
	}

	fmt.Println("read : ", done)
	c.L.Unlock()
}

func write(c *sync.Cond) {
	c.L.Lock()
	time.Sleep(time.Second)
	fmt.Println("write func signal")
	// 写入完成标记
	done = true
	c.L.Unlock()
	// 唤醒等待的资源
	c.Broadcast()
}
