package test_best_practices

import "sync"

type gracefulTerminationManager struct {
	rsList graceTerminateRSList
}

func newGracefulTerminationManager() *gracefulTerminationManager {
	return &gracefulTerminationManager{
		rsList: graceTerminateRSList{
			list: make(map[string]*item),
		},
	}
}

type item struct {
	VirtualServer string
	RealServer    string
}

type graceTerminateRSList struct {
	lock sync.Mutex
	list map[string]*item
}

func (g *graceTerminateRSList) flushList(handler func(rsToDelete *item) (bool, error)) bool {
	g.lock.Lock()
	defer g.lock.Unlock()
	success := true
	for _, rs := range g.list {
		if ok, err := handler(rs); !ok || err != nil {
			success = false
		}
	}
	return success
}

func (g *graceTerminateRSList) add(rs *item) {
	g.lock.Lock()
	defer g.lock.Unlock()
	g.list[rs.RealServer] = rs
}

func (g *graceTerminateRSList) len() int {
	g.lock.Lock()
	defer g.lock.Unlock()
	return len(g.list)
}
