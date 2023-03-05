package define

import "sync"

type SegmentCache struct {
	cache            sync.Map
	loadSegmentEvent chan SegmentType // 加载号段的信号量
}

func NewSegmentCache() *SegmentCache {
	segmentCache := &SegmentCache{
		loadSegmentEvent: make(chan SegmentType, 8),
	}
	return segmentCache
}

func (s *SegmentCache) Add(alloc *SegmentAlloc) SegmentType {
	s.cache.Store(alloc.SegmentType, alloc)
	return alloc.SegmentType
}

func (s *SegmentCache) Get(segmentType SegmentType) *SegmentAlloc {
	v, ok := s.cache.Load(segmentType)
	if ok {
		return v.(*SegmentAlloc)
	}
	return nil
}

func (s *SegmentCache) LoadEvent() <-chan SegmentType {
	return s.loadSegmentEvent
}

func (s *SegmentCache) WriteLoadEvent(segmentAlloc *SegmentAlloc) {
	segmentAlloc.IsPreloading = true
	// 这里有可能缓冲区被写满进入等待，缓冲区目前为8
	s.loadSegmentEvent <- segmentAlloc.SegmentType
}
