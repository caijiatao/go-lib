package impl

import (
	"sort"
	"sync"
	"testing"
	"time"

	mapset "github.com/deckarep/golang-set"
	"github.com/stretchr/testify/assert"

	"golib/id_generate/internal/define"
	"golib/id_generate/service/impl/snowflake_generator"
)

func TestGetSegmentId(t *testing.T) {
	var (
		segmentType = define.MockBusinessSegmentType
		err         error
	)
	// 初始化完成
	for NewIdGenerateService(snowflake_generator.NewSnowflakeGenerateService()).segmentCache.Get(segmentType) == nil {
		time.Sleep(time.Second)
	}
	step := testIdGenerateService.segmentCache.Get(segmentType).Step
	var (
		generateIds     = make([]uint64, 0)
		generateIdMutex sync.Mutex
	)
	// 重新初始化一下cache，避免号段被污染
	err = testIdGenerateService.initCacheAllSegmentAlloc(testCtx)
	assert.Nil(t, err)
	var wg sync.WaitGroup

	// 1.模拟客户端并发获取id号
	var i uint64
	for i = uint64(0); i < 2*step; i++ {
		wg.Add(1)
		go func() {
			nextId, err := testIdGenerateService.GetSegmentId(testCtx, segmentType)
			assert.Nil(t, err)
			generateIdMutex.Lock()
			defer generateIdMutex.Unlock()
			generateIds = append(generateIds, nextId)
			wg.Done()
		}()
	}
	wg.Wait()

	// 2.判断所有id不重号
	assert.Equal(t, i, uint64(len(generateIds)))
	generateIdSet := mapset.NewSet()
	for _, v := range generateIds {
		assert.False(t, generateIdSet.Contains(v))
		generateIdSet.Add(v)
	}

	// 3.判断是否递增
	generateIdsUintSlice := generateIdSet.ToSlice()
	sort.Slice(generateIdsUintSlice, func(i, j int) bool {
		return generateIdsUintSlice[i].(uint64) < generateIdsUintSlice[j].(uint64)
	})
	for i := range generateIdsUintSlice {
		if i == 0 {
			continue
		}
		assert.Greater(t, generateIdsUintSlice[i], generateIdsUintSlice[i-1])
	}
}
