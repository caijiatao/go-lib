package impl

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"

	internalDefine "golib/id_generate/internal/define"
	"golib/id_generate/internal/repository"
	"golib/id_generate/service/impl/snowflake_generator"
)

var (
	idGenerateService         *IdGenerateService
	idGenerateServiceInitOnce sync.Once
)

func NewIdGenerateService(snowflakeGenerator *snowflake_generator.Node) *IdGenerateService {
	idGenerateServiceInitOnce.Do(func() {
		var (
			ctx, cancel = context.WithCancel(context.Background())
		)
		idGenerateService = &IdGenerateService{
			segmentAPI:         repository.NewSegmentAPI(),
			snowflakeGenerator: snowflakeGenerator,
			segmentCache:       internalDefine.NewSegmentCache(),
		}

		go idGenerateService.initSegmentCache(ctx)

		// 优雅退出
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		go func() {
			osCall := <-c
			log.Printf("system call: %+v", osCall)
			cancel()
		}()
	})
	return idGenerateService
}
