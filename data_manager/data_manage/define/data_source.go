package define

// SourcePGDBConfig
// @Description: 源数据库信息
type SourcePGDBConfig struct {
	DataSourceId uint64
	Port         string
	Host         string
	User         string
	Pwd          string // 密码加密存储
	DbName       string
}

type DataSourceConfig struct {
	SourceType uint64
	SourcePGDBConfig
}

func NewDataSourceConfig(sourceType uint64) *DataSourceConfig {
	return &DataSourceConfig{SourceType: sourceType}
}
