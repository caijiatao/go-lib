package lruk

import (
	lru "github.com/hashicorp/golang-lru/v2"
	"sync"
	"time"
)

type Key interface {
	Len() int
}

type Value interface {
	Len() int
}

type entry struct {
	value      Value
	k          int
	lastAccess time.Time

	mu sync.Mutex
}

func (e *entry) IncrementAndCheckK(k int) bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.k++
	e.lastAccess = time.Now()
	return e.k >= k
}

type LRUK struct {
	cache *lru.Cache[Key, *entry]

	config
}

type config struct {
	size int
	k    int
}

type Opt func(*config)

func WithK(k int) Opt {
	return func(l *config) {
		l.k = k
	}
}

func WithSize(size int) Opt {
	return func(l *config) {
		l.size = size
	}
}

func NewLRUK(opts ...Opt) (*LRUK, error) {
	config := config{
		size: 10000,
		k:    3,
	}
	for _, opt := range opts {
		opt(&config)
	}

	baseCache, err := lru.New[Key, *entry](config.size)
	if err != nil {
		return nil, err
	}
	return &LRUK{
		cache:  baseCache,
		config: config,
	}, nil
}

func (l *LRUK) Get(key Key) (value Value, ok bool) {
	e, ok := l.cache.Get(key)
	if !ok {
		return
	}

	if e.IncrementAndCheckK(l.k) {
		l.cache.Add(key, e)
	}

	value = e.value

	return
}

func (l *LRUK) Add(key Key, value Value) {
	l.cache.Add(
		key,
		&entry{
			value:      value,
			lastAccess: time.Now(),
		},
	)
}

func (l *LRUK) Remove(key Key) {
	l.cache.Remove(key)
}
