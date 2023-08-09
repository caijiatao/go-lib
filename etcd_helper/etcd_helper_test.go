package etcd_helper

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPing(t *testing.T) {
	InitTestSuite()

	ctx := BindContext(context.Background())
	err := Context(ctx).Ping(ctx)
	assert.Nil(t, err)
}
