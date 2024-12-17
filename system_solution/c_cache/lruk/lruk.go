package lruk

import (
	lru "github.com/hashicorp/golang-lru/v2"
	"sync"
)

type Key interface {
	Len() int
}

type Value interface {
	Len() int
}

type entry struct {
	// 具体的值
	value Value
	// 访问次数
	k int

	mu sync.Mutex
}

func (e *entry) IncrementAndCheckK(k int) bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.k++
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
	c := config{
		size: 10000,
		k:    3,
	}
	for _, opt := range opts {
		opt(&c)
	}

	baseCache, err := lru.New[Key, *entry](c.size)
	if err != nil {
		return nil, err
	}
	return &LRUK{
		cache:  baseCache,
		config: c,
	}, nil
}

func (l *LRUK) Get(key Key) (value Value, ok bool) {
	e, ok := l.cache.Get(key)
	if !ok {
		return
	}

	// 如果访问次数达到了k次，那么就将其放到队头，避免被淘汰
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
			value: value,
		},
	)
}

func (l *LRUK) Remove(key Key) {
	l.cache.Remove(key)
}
