package repository

import (
	"context"
	"golib/libs/id_generate/internal/define"
	"sync"
)

type mockIdSegmentImpl struct {
	sync.Mutex
	localCache map[define.SegmentType]uint64
}

func newMockIdSegmentImpl() *mockIdSegmentImpl {
	return &mockIdSegmentImpl{
		localCache: map[define.SegmentType]uint64{},
	}
}

func (m *mockIdSegmentImpl) GetNextIdSegment(ctx context.Context, segmentType define.SegmentType) (start uint64, end uint64, err error) {
	m.Lock()
	if _, ok := m.localCache[segmentType]; ok {
		start = m.localCache[segmentType]
		m.localCache[segmentType] = start + 100
	} else {
		start = 0
		m.localCache[segmentType] = 100
	}
	end = m.localCache[segmentType]
	m.Unlock()
	return start, end, nil
}
