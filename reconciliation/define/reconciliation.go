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
