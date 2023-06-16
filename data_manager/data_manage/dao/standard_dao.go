package dao

import (
	"airec_server/pkg/gorm_helper"
	"context"
)

func (dao *DataManageDao) SaveStandardData(ctx context.Context, tableName string, data interface{}) error {
	gorm_helper.Context(ctx).Table(tableName).Save(data)
	return nil
}
