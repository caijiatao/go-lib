package repository

import (
	"context"
	"golib/libs/id_generate/internal/define"
)

type SegmentAPI interface {
	GetNextIdSegment(ctx context.Context, segmentType define.SegmentType) (start, end uint64, err error)
}
