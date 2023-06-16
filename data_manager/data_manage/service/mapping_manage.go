package service

import (
	"context"
)

type MappingManage struct{}

func (m *MappingManage) RunSyncDataMapping(ctx context.Context, value syncDataJobValue) {
	// TODO 并发数控制
	mappingProcessor, err := NewMappingProcessor(ctx, value)
	if err != nil {
		return
	}

	mappingProcessor.Run(ctx)
}
