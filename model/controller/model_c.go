package modelController

type ReturnIsNotFoundOnly struct {
	IsNotFound bool
}

type ParamPagination struct {
	Page        int64
	Limit       int64
	SortField   string
	SortOrderBy string
}

type ReturnPagination struct {
	TotalData int64
	TotalPage int64
	Page      int64
	Limit     int64
}