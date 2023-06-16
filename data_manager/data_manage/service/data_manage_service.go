package service

import (
	"airec_server/internal/data_manage/dao"
	"airec_server/internal/data_manage/define"
	"context"
)

type IDataManageService interface {
	GetSourceMappingBySourceId(ctx context.Context, dataSourceId uint64) (*define.SourceMappingDto, error)

	CreateDataSource(ctx context.Context, config define.DataSourceConfig) (err error)
	GetDataSourceList(ctx context.Context, req define.GetDataSourceListReq) (resp define.GetDataSourceListResp, err error)
	GetDataSourceDetails(ctx context.Context, req define.GetDataSourceDetailsReq) (resp define.GetDataSourceDetailsResp, err error)
}

type DataManageService struct {
	IDataManageService
	dataManageDao *dao.DataManageDao
}

func GetDataManageService() IDataManageService {
	return &DataManageService{}
}

func (dms *DataManageService) GetSourceMappingBySourceId(ctx context.Context, dataSourceId uint64) (*define.SourceMappingDto, error) {
	sourceModel, err := dao.GetDataManageDao().GetSourceByDatasourceId(ctx, dataSourceId)
	if err != nil {
		return nil, err
	}
	sourceDto := &define.SourceMappingDto{
		SourceModel: sourceModel,
	}
	return sourceDto, nil
}
