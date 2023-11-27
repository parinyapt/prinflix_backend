package modelHandler

type QueryParamPagination struct {
	Page        int64  `form:"page" validate:"omitempty,min=1"`
	Limit       int64  `form:"limit" validate:"omitempty,min=1,max=100"`
	SortField   string `form:"sortby" validate:"omitempty"`
	SortOrderBy string `form:"orderby" validate:"omitempty,oneof=asc desc"`
}

type ResponsePagination struct {
	TotalData    int64 `json:"total_data"`
	TotalPage    int64 `json:"total_page"`
	CurrentPage  int64 `json:"current_page"`
	CurrentLimit int64 `json:"current_limit"`
}