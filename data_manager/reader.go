package data_manager

import "context"

type SourceDataReader interface {
	Read(ctx context.Context) ([]map[string]interface{}, error)
}
