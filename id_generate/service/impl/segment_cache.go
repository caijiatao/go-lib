package impl

import (
	"context"
	"fmt"
	"log"

	"github.com/pkg/errors"

	"golib/goasync"
	"golib/id_generate/internal/define"
)

func (service *IdGenerateService) initSegmentCache(ctx context.Context) {
	log.Println(ctx, "init segment cache start")
	// 不强绑定在服务启动
	defer func() {
		r := recover()
		_ = goasync.PanicErrHandler(r)
	}()
	// 初始化失败直接返回
	err := idGenerateService.initCacheAllSegmentAlloc(ctx)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Printf("init segment cache success\n")
	log.Println(ctx, "init segment cache success")
	idGenerateService.watchSegmentLoadEvent(ctx)
}

func (service *IdGenerateService) initCacheAllSegmentAlloc(ctx context.Context) error {
	if err := service.initCacheSegmentAlloc(ctx, define.MockBusinessSegmentType, 100); err != nil {
		return err
	}
	return nil
}

func (service *IdGenerateService) initCacheSegmentAlloc(ctx context.Context, segmentType define.SegmentType, step uint64) error {
	segmentAlloc := define.NewSegmentAlloc(segmentType, step)
	segmentAlloc.IsPreloading = true
	err := service.loadSegmentAllocBuffer(ctx, segmentAlloc)
	if err != nil {
		return errors.WithStack(err)
	}
	service.segmentCache.Add(segmentAlloc)
	return nil
}

func (service *IdGenerateService) watchSegmentLoadEvent(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Printf("watchSegmentLoadEvent stop")
			return
		case segmentType, ok := <-service.segmentCache.LoadEvent():
			if !ok {
				continue
			}
			err := service.preloadSegmentAllocBufferBySegmentType(ctx, segmentType)
			if err != nil {
				log.Printf("loadSegmentAllocBufferErr:%#v", err)
				continue
			}
		}
	}
}

// preloadSegmentAllocBufferBySegmentType
// @Description: 通过号段类型预先加载号段
func (service *IdGenerateService) preloadSegmentAllocBufferBySegmentType(ctx context.Context, segmentType define.SegmentType) (err error) {
	segmentAlloc := service.segmentCache.Get(segmentType)
	if segmentAlloc == nil {
		return nil
	}
	// 修改preloading时上锁
	segmentAlloc.Lock()
	defer func() {
		segmentAlloc.IsPreloading = false
		segmentAlloc.Unlock()
		if err == nil {
			wakeupErr := segmentAlloc.WakeUpAllWaitingClient()
			if wakeupErr != nil {
				log.Printf("wake up all client err : %+v", wakeupErr)
			}
		}
	}()
	err = service.loadSegmentAllocBuffer(ctx, segmentAlloc)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// loadSegmentAllocBuffer
// @Description: 加载id号，使用时直接从内存中获取，这里不会对alloc进行加锁，
//	如果是预加载则需要加锁，非预加载的情况不会被并发访问，因为还没有加到segment cache 里面
func (service *IdGenerateService) loadSegmentAllocBuffer(ctx context.Context, segmentAlloc *define.SegmentAlloc) (err error) {
	defer func() {
		r := recover()
		_ = goasync.PanicErrHandler(r)
	}()
	segment := define.NewSegment(segmentAlloc.Step)
	log.Printf("idCenterApi loadSegmentAllocBuffer: %+v", segmentAlloc)
	for segmentAlloc.Step >= uint64(len(segment.Ids)) {
		start, end, err := service.segmentAPI.GetNextIdSegment(ctx, segmentAlloc.SegmentType)
		if err != nil {
			log.Printf("idCenterApi.GetSeqIdBySegmentTypeErr :%+v", err)
			continue
		}
		for i := start; i < end; i++ {
			segment.Ids = append(segment.Ids, i)
		}
	}
	segmentAlloc.SegmentBuffers = append(segmentAlloc.SegmentBuffers, segment)
	log.Printf("idCenterApi loadSegmentAllocBuffer finish: %+v", segmentAlloc)
	return nil
}
