package orm

const (
	TimeSortDefault = iota
	TimeSortDesc
	TimeSortAsc
)

type OrderByParams struct {
	TimeSort uint8 `json:"timeSort" form:"timeSort"`
}

func (o *OrderByParams) GetOrderBys() []string {
	switch o.TimeSort {
	case TimeSortDesc:
		return []string{"id DESC"}
	case TimeSortAsc:
		return []string{"id ASC"}
	}
	return nil
}
