package define

type ReconciliationArgs struct{}

type ReconciliationResult struct {
	CreateTime   uint
	BusinessType string
	BatchNum     string

	Total    uint
	Succ     uint
	SuccRate float32
	Fail     uint
	Remark   string
}

type DatabaseConfig struct {
	Addr     string
	Port     string
	User     string
	Password string
	TableDB  string
	Table    string
}

type FieldsConfig struct {
	// 原始字段映射到目标字段，如果没有名字变更则无需特别设置
	FiledMapping map[string]string
	// 唯一字段，通过该字段来进行查询
	SourceUniqFields []string
	// 不需要对比的字段，eg:如果进行了分库，那么id大概率是不需要对比的，可以直接进行忽略
	IgnoreFields []string
}
