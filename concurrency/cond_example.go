package concurrency

import (
	"fmt"
	"sync"
	"time"
)

var done bool

func read(c *sync.Cond) {
	c.L.Lock()

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
	done = true
	c.L.Unlock()
	c.Broadcast()
}
