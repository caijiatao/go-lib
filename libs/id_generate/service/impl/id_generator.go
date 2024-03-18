package impl

import (
	"context"
	"golib/libs/id_generate/define"
	define2 "golib/libs/id_generate/internal/define"
	"log"
	"time"

	"github.com/pkg/errors"
)

func (service *IdGenerateService) GetSegmentId(ctx context.Context, segmentType define2.SegmentType) (id uint64, err error) {
	segmentAlloc := service.segmentCache.Get(segmentType)
	defer func() {
		if err == nil {
			return
		}
		id, _, err = service.segmentAPI.GetNextIdSegment(ctx, segmentType)
		if err != nil {
			return
		}
	}()
	if segmentAlloc == nil {
		// 还没有预分配好内存，直接从数据库获取
		err = define.RespSegmentAllocNotInitErr
		return 0, err
	}

	id, err = service.nextId(ctx, segmentAlloc)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	return id, nil
}

func (service *IdGenerateService) nextId(ctx context.Context, segmentAlloc *define2.SegmentAlloc) (id uint64, err error) {
	if segmentAlloc == nil {
		log.Println(ctx)
		return 0, define.RespGenerateIdErr
	}

	segmentAlloc.Lock()
	defer segmentAlloc.Unlock()

	id = segmentAlloc.GetId()

	if segmentAlloc.IsNeedPreload() {
		service.segmentCache.WriteLoadEvent(segmentAlloc)
	}

	// 如果已经获取到正确的id则直接返回
	if segmentAlloc.IsRightId(id) {
		return id, nil
	}

	// 如果没拿到号段 ，在这里加入等待队列，前面已经发出事件开始加载，避免多个协程同时进行加载
	waitChan := make(chan struct{}, 1)
	segmentAlloc.Waiting = append(segmentAlloc.Waiting, waitChan)
	// 让其他客户端可以走前面的步骤，进入到等待状态
	segmentAlloc.Unlock()

	// 最多等待500ms，超过等待时间则直接返回错误
	timer := time.NewTimer(50 * time.Millisecond)
	select {
	case <-waitChan:
	case <-timer.C:
	}

	segmentAlloc.Lock()
	segmentAlloc.RefreshBuffer()
	id = segmentAlloc.GetId()
	if segmentAlloc.IsRightId(id) {
		return id, nil
	}
	log.Println(ctx)
	return 0, define.RespGenerateIdErr
}
