package concurrency

import (
	"fmt"
	"time"
)

func WatchLeases() {
	leaseWatchChan := make(chan []int)

	go func() {
		count := 0
		for {
			time.Sleep(time.Second)
			count++
			leaseWatchChan <- []int{count}
		}
	}()

	for n := range leaseWatchChan {
		fmt.Println(n)
	}
}
