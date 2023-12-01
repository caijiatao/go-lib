package concurrency

import (
	"fmt"
	"sync"
	"time"
)

func lockRead(mutex *sync.Mutex) {
	mutex.Lock()

	fmt.Println("read")
	time.Sleep(time.Second)

	mutex.Unlock()
}

func lockWrite(mutex *sync.Mutex) {
	mutex.Lock()

	time.Sleep(time.Second)
	fmt.Println("write")

	mutex.Unlock()
}
