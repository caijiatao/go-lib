package orm

type PageParams struct {
	PageSize int `json:"pageSize" binding:"required" form:"pageSize"`
	PageNum  int `json:"pageNum"  binding:"required" form:"pageNum"`
}

type PageCommonResponse struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	PageNum  int         `json:"pageNum"`
	PageSize int         `json:"pageSize"`
}

func NewPageCommonResponse(list interface{}, total int64, pageParams PageParams) *PageCommonResponse {
	return &PageCommonResponse{
		List:     list,
		Total:    total,
		PageNum:  pageParams.PageNum,
		PageSize: pageParams.PageSize,
	}
}
