package impl

import (
	"context"

	"golib/id_generate/internal/define"
	"golib/id_generate/internal/repository"
	"golib/id_generate/service/impl/snowflake_generator"
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
