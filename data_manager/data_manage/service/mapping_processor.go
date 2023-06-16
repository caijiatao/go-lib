package service

import (
	"airec_server/internal/data_manage/define"
	"airec_server/pkg/goasync"
	"context"
	"errors"
	"sync"
)

type MappingProcessor struct {
	// 数据源获取，每次映射只支持一个数据源
	reader IDataReader
	// 原始数据写入
	rawDataWriters []DataWriter
	// 映射处理完成后的数据写入
	mappedDataWriters []DataWriter

	// 源表映射到标准表
	sourceTable2StandardTableMappingMap map[string]StandardTableMapping
}

func NewMappingProcessor(ctx context.Context, jobValue syncDataJobValue) (*MappingProcessor, error) {
	mappingProcessor := &MappingProcessor{
		sourceTable2StandardTableMappingMap: make(map[string]StandardTableMapping),
	}
	sourceMappingDto, err := GetDataManageService().GetSourceMappingBySourceId(ctx, jobValue.DatasourceId)
	if err != nil {
		return nil, err
	}

	// 将数据源的所有表都映射到标准模型上
	for sourceTableId, sourceTableModel := range sourceMappingDto.SourceTableId2SourceTableMap {
		standardTableModel, ok := sourceMappingDto.SourceTableId2StandardTableMap[sourceTableId]
		if !ok {
			continue
		}
		standardTableMapping := StandardTableMapping{
			MappingTableName:            standardTableModel.TableName,
			ColumnName2ColumnMappingMap: make(map[string]ColumnStandardMapping),
		}
		standardMappingDto := sourceMappingDto.GetColumnInfo(sourceTableId, sourceTableId)
		standardTableMapping.ColumnName2ColumnMappingMap[standardMappingDto.SourceColumn.ColumnName] =
			NewStandardMapping(standardMappingDto.GetStandardColumnType(), standardMappingDto.GetConvertScript(), standardMappingDto.GetStandardColumnName())

		mappingProcessor.sourceTable2StandardTableMappingMap[sourceTableModel.TableName] = standardTableMapping
	}

	mappingProcessor.reader = NewDataReader(define.NewDataSourceConfig(sourceMappingDto.SourceModel.SourceType))
	mappingProcessor.rawDataWriters = newRawDataWriter()
	mappingProcessor.mappedDataWriters = newDefaultMappingWriter()
	return mappingProcessor, nil
}

func newRawDataWriter() []DataWriter {
	//TODO 原始数据保存，这里后续可以做成在配置数据源的时候就进行配置
	return nil
}

func newDefaultMappingWriter() []DataWriter {
	return nil
}

func (m *MappingProcessor) Run(ctx context.Context) {
	var rawDataChan chan Data
	// 读取原始数据
	go func() {
		defer func() {
			r := recover()
			_ = goasync.PanicErrHandler(r)
		}()
		rawDataChan = m.reader.Read(ctx)
	}()

	for _, writer := range m.mappedDataWriters {
		writer.Run(ctx)
	}

	// 原始数据广播给Writer进行清洗及存储
	go func() {
		defer func() {
			r := recover()
			_ = goasync.PanicErrHandler(r)
		}()
		for {
			select {
			case <-ctx.Done():
				return
			case data, ok := <-rawDataChan:
				if !ok {
					return
				}
				m.broadcast(data)
			}
		}
	}()
}

// broadcast
//
//	@Description: 数据广播给所有Writer
func (m *MappingProcessor) broadcast(data Data) {
	// 原始数据广播
	var wg sync.WaitGroup
	wg.Add(len(m.rawDataWriters) + len(m.mappedDataWriters))
	for _, writer := range m.rawDataWriters {
		go func(dw DataWriter, d Data) {
			defer func() {
				r := recover()
				_ = goasync.PanicErrHandler(r)
			}()
			dw.DataChan <- d
			wg.Done()
		}(writer, data)
	}

	mappedData, err := m.mappingToStandardData(data)
	if err != nil {
		// TODO handle error
	}

	for _, writer := range m.mappedDataWriters {
		go func(dw DataWriter, d Data) {
			defer func() {
				r := recover()
				_ = goasync.PanicErrHandler(r)
			}()
			dw.DataChan <- d
			wg.Done()
		}(writer, mappedData)
	}

	// 等待所有都广播完成才处理下一条数据
	wg.Wait()
}

func (m *MappingProcessor) Close() error {
	//  TODO close all resource
	return nil
}

// mappingToStandardData
//
//	@Description: 处理数据映射，转化成标准模型
//	@return standDataResults 标准化后的数据结果
func (m *MappingProcessor) mappingToStandardData(data Data) (mappedStandardData Data, err error) {
	var (
		ok                   bool
		dataMap              map[string]interface{}
		standardTableMapping StandardTableMapping
	)
	if dataMap, ok = data.RawData.(map[string]interface{}); !ok {
		return mappedStandardData, errors.New("raw data type error")
	}
	if standardTableMapping, ok = m.sourceTable2StandardTableMappingMap[data.TableName]; !ok {
		return mappedStandardData, errors.New("mapping error")
	}
	standardDataResult := make(map[string]interface{})
	for columnName, value := range dataMap {
		// 找到对应的映射信息
		var columnMapping ColumnStandardMapping
		if columnMapping, ok = standardTableMapping.ColumnName2ColumnMappingMap[columnName]; !ok {
			// 沒有映射
			continue
		}

		// 映射处理成标准值
		v, err := columnMapping.MapToStandValue(value)
		if err != nil {
			return mappedStandardData, err
		}
		standardDataResult[columnMapping.StandardColumnName] = v
	}

	// 赋值映射好后的数据
	mappedStandardData.TableName = standardTableMapping.MappingTableName
	mappedStandardData.RawData = standardDataResult
	return mappedStandardData, nil
}
