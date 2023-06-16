package define

import "airec_server/internal/data_manage/model"

type SourceMappingDto struct {
	SourceModel                  *model.Source
	SourceTableId2SourceTableMap map[uint64]model.SourceTable
	// sourceTableId & sourceColumnId ->  StandardMappingDto
	Source2StandardColumnMap        map[uint64]map[uint64]StandardMappingDto
	SourceTableId2SourceColumnIdMap map[uint64][]uint64
	SourceColumnId2SourceColumnMap  map[uint64]model.SourceColumn
	SourceTableId2StandardTableMap  map[uint64]model.StandardTable

	StandardMappingModel []model.StandardMapping

	StandardTableId2StandardTableMap map[uint64]model.StandardTable
}

type StandardMappingDto struct {
	model.SourceColumn
	model.StandardTableColumn
	model.StandardMapping
}

func (dto *StandardMappingDto) GetStandardColumnType() string {
	return dto.StandardTableColumn.ColumnType
}

func (dto *StandardMappingDto) GetConvertScript() string {
	return dto.StandardMapping.ConvertScript
}

func (dto *StandardMappingDto) GetStandardColumnName() string {
	return dto.StandardTableColumn.ColumnName
}

func (dto *SourceMappingDto) GetSourceTableColumnIdsByTableId(tableId uint64) []uint64 {
	return dto.SourceTableId2SourceColumnIdMap[tableId]
}

func (dto *SourceMappingDto) GetColumnInfo(sourceTableId uint64, sourceColumnId uint64) (standardMappingDto StandardMappingDto) {
	if standardColumnMap, ok := dto.Source2StandardColumnMap[sourceTableId]; ok {
		if standardMappingDto, ok = standardColumnMap[sourceColumnId]; ok {
			return standardMappingDto
		}
	}
	return standardMappingDto
}
