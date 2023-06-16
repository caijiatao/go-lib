package define

import "airec_server/internal/data_manage/model"

type GetDataSourceListReq struct {
}

type GetDataSourceListResp struct {
	model.Source
	DataSourceConfig
}

type GetDataSourceDetailsReq struct {
}

type GetDataSourceDetailsResp struct {
}
