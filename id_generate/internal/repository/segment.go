package repository

import (
	"context"

	"golib/id_generate/internal/define"
)

type idSegmentRepository struct {
}

func (repo *idSegmentRepository) GetNextIdSegment(ctx context.Context, segmentType define.SegmentType) (uint64, uint64, error) {
	// implement me , get id segment from DB ，每个DB实现不同，这里不做具体实现
	return 0, 0, nil
}

func newIdSegmentRepository() *idSegmentRepository {
	return &idSegmentRepository{}
}
