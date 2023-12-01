package concurrency

import (
	"sync"
	"testing"
	"time"
)

func Test_write(t *testing.T) {
	cond := sync.NewCond(&sync.Mutex{})
	go read(cond)

	go read(cond)

	time.Sleep(time.Second * 1)

	write(cond)

	time.Sleep(time.Second * 3)
}
