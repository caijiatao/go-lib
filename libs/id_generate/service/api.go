package service

import "context"

type IdGenerateApi interface {
	GenerateId(ctx context.Context) (string, error)
}
