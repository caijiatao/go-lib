package concurrency

import (
	"sync"
	"testing"
	"time"
)

func Test_lockWrite(t *testing.T) {
	mutex := &sync.Mutex{}
	go lockRead(mutex)

	go lockRead(mutex)

	time.Sleep(time.Second)
	lockWrite(mutex)
	time.Sleep(time.Second * 3)
}
