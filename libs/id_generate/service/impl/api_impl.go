package impl

import (
	"context"
	"golib/libs/id_generate/internal/define"
	"golib/libs/id_generate/internal/repository"
	"golib/libs/id_generate/service/impl/snowflake_generator"
)

type IdGenerateService struct {
	snowflakeGenerator *snowflake_generator.Node
	segmentAPI         repository.SegmentAPI
	segmentCache       *define.SegmentCache
}

func (service *IdGenerateService) GenerateSnowflakeId(ctx context.Context) (uint64, error) {
	id := service.snowflakeGenerator.Generate().Uint64()
	return id, nil
}
