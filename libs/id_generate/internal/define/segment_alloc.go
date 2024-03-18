package define

import (
	"golib/libs/goasync"
	"sync"
	"sync/atomic"
	"time"
)

type SegmentAlloc struct {
	SegmentType               SegmentType     // 号段业务类型
	CurrentSegmentBufferIndex int             // 当前使用的buffer 下标，会一直使用0的下标
	SegmentBuffers            []*Segment      // 号段的双buffer
	Step                      uint64          // 号段加载的步长
	IsPreloading              bool            // 是否处于预加载号段的状态
	UpdateTime                time.Time       // 更新号段的时间
	mutex                     sync.Mutex      // 互斥锁
	Waiting                   []chan struct{} //等待的客户端，当号段用户正处于初始化时，其他协程处于等待状态，到一定超时时间仍然未完成再返回失败
}

func NewSegmentAlloc(segmentType SegmentType, step uint64) *SegmentAlloc {
	return &SegmentAlloc{
		SegmentType:    segmentType,
		Step:           step,
		UpdateTime:     time.Now(),
		SegmentBuffers: make([]*Segment, 0),
		Waiting:        make([]chan struct{}, 0),
	}
}

func (s *SegmentAlloc) Lock() {
	s.mutex.Lock()
}

func (s *SegmentAlloc) Unlock() {
	s.mutex.Unlock()
}

func (s *SegmentAlloc) IsHasSegment() bool {
	currentBuffer := s.SegmentBuffers[s.CurrentSegmentBufferIndex]
	// 游标已经到最大游标则没有可取的id了
	if currentBuffer.Cursor <= currentBuffer.MaxCursor {
		return true
	}
	return false
}

func (s *SegmentAlloc) GetId() uint64 {
	if s.IsHasSegment() {
		currentBuffer := s.SegmentBuffers[s.CurrentSegmentBufferIndex]
		// 获取新的id，并将下标后移一位
		id := currentBuffer.Ids[currentBuffer.Cursor]
		atomic.AddUint64(&currentBuffer.Cursor, 1)
		s.UpdateTime = time.Now()
		return id
	}
	return 0
}

func (s *SegmentAlloc) IsRightId(id uint64) bool {
	return id > 0
}

func (s *SegmentAlloc) IsNeedPreload() bool {
	// 已经在预加载了
	if s.IsPreloading {
		return false
	}
	// 第二个缓冲区已经准备好 ，这里之前遗漏了该判断，会导致只要超过一半就开始去预加载
	if len(s.SegmentBuffers) > 1 {
		return false
	}
	segmentBuffer := s.SegmentBuffers[s.CurrentSegmentBufferIndex]
	// 当前剩余的号已经小于步长的一半，则进行加载
	restId := segmentBuffer.MaxCursor - segmentBuffer.Cursor
	if restId <= s.Step/2 {
		return true
	}
	return false
}

func (s *SegmentAlloc) WakeUpAllWaitingClient() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, waiting := range s.Waiting {
		goasync.SafeClose(waiting)
	}
	return nil
}

func (s *SegmentAlloc) IsNewBufferReady() bool {
	if len(s.SegmentBuffers) <= 1 {
		return false
	}
	return true
}

func (s *SegmentAlloc) RefreshBuffer() {
	// 当前buffer 仍然有号则不需要刷新，因为可能其他协程已经刷新了buffer区
	if s.IsHasSegment() {
		return
	}
	if !s.IsNewBufferReady() {
		return
	}
	s.SegmentBuffers = append(s.SegmentBuffers[:0], s.SegmentBuffers[1:]...)
}
