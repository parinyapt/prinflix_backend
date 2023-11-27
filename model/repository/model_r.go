package modelRepository

type ResultIsFoundOnly struct {
	IsFound bool
}

type ParamPagination struct {
	Page        int64
	Limit       int64
	SortField   string
	SortOrderBy string
}

type ResultPagination struct {
	TotalData int64
	TotalPage int64
	Page      int64
	Limit     int64
}
