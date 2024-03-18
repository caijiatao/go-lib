package context_helper

import (
	"airec-libs/etcd_helper"
	"airec-libs/orm"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetKeyValues(t *testing.T) {
	orm.InitTestSuite()
	ctx := orm.BindContext(context.Background())

	etcd_helper.InitTestSuite()
	ctx = etcd_helper.BindContext(ctx)
	keyValues := GetKeyValues(ctx)

	assert.Equal(t, 2, len(keyValues))
}
