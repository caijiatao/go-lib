package dao

import (
	"airec_server/internal/data_manage/model"
	"airec_server/pkg/gorm_helper"
	"context"
)

func (dao *DataManageDao) GetSourceByDatasourceId(ctx context.Context, datasourceId uint64) (*model.Source, error) {
	source := &model.Source{}
	gorm_helper.Context(ctx).Model(model.Source{}).Where("datasource_id = ?", datasourceId).First(source)
	return source, nil
}
