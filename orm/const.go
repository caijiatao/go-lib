package orm

const (
	PGSQLSourceType = iota + 1
	FileSourceType
	MysqlSourceType
	HiveSourceType
	HDFSSourceType
	OracleSourceType
	KafkaSourceType
)

const (
	DefaultBatchCreateSize = 5000
)
