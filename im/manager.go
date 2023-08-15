package im

import "sync"

type Manager struct {
	clients sync.Map
}
