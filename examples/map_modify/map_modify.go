package map_modify

import (
	"sync"
)

const mapSize = 1000
const numIterations = 100000

func writeToMapWithMutex() {
	m := make(map[int]int)
	var mutex sync.Mutex

	for i := 0; i < numIterations; i++ {
		mutex.Lock()
		m[i%mapSize] = i
		mutex.Unlock()
	}
}

func writeToMapWithChannel() {
	m := make(map[int]int)
	ch := make(chan struct {
		key   int
		value int
	}, 256)

	var wg sync.WaitGroup
	go func() {
		wg.Add(1)
		for {
			entry, ok := <-ch
			if !ok {
				wg.Done()
				return
			}
			m[entry.key] = entry.value
		}
	}()

	for i := 0; i < numIterations; i++ {
		ch <- struct {
			key   int
			value int
		}{i % mapSize, i}
	}
	close(ch)
	wg.Wait()
}
