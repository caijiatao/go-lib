package snowflake_generator

import (
	"sync"
	"testing"
)

const (
	nodeCnt = 10
	idCnt   = 10000
)

func TestNodeGenerate(t *testing.T) {
	var wg sync.WaitGroup
	channels := make([]chan uint64, nodeCnt, nodeCnt)
	Nodes := make([]*Node, nodeCnt, nodeCnt)
	for i := 0; i < nodeCnt; i++ {
		channels[i] = make(chan uint64, idCnt)
		Nodes[i] = NewSnowflakeGenerateService()
	}
	for i := 0; i < nodeCnt; i++ {
		wg.Add(1)
		channel := channels[i]
		Node := Nodes[i]
		go func() {
			defer wg.Done()
			for j := 0; j < idCnt; j++ {
				channel <- Node.Generate().Uint64()
			}
		}()
	}
	wg.Wait()
	idSet := make(map[uint64]bool)
	for i := range channels {
		for id := range channels[i] {
			if idSet[id] {
				t.Error("Duplicate ID")
				return
			}
			idSet[id] = true
			if len(channels[i]) == 0 {
				break
			}
		}
	}
}
