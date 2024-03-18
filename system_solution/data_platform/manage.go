package data_platform

import "context"

type Manager interface {
	WatchData(ctx context.Context, receiver []IncrementData) error
}
