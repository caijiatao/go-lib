package service

import (
	"airec_server/internal/data_manage/define"
	"context"
)

type ReaderConstructor func(config *define.DataSourceConfig) IDataReader

var (
	dataReaderConstructorMap = map[uint64]ReaderConstructor{
		define.PGSQLSourceType: NewPGSQLDataReader,
	}
)

type IDataReader interface {
	Read(ctx context.Context) chan Data
}

func NewDataReader(config *define.DataSourceConfig) IDataReader {
	return dataReaderConstructorMap[config.SourceType](config)
}

type PGSQLDataReader struct{}

func NewPGSQLDataReader(config *define.DataSourceConfig) IDataReader {
	// TODO 初始化数据库
	return nil
}
