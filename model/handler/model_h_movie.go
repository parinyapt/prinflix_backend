package modelHandler

type ResponseGetMovieCategoryList struct {
	CategoryID   uint   `json:"id"`
	CategoryName string `json:"name"`
}

type QueryParamPagination struct {
	Page  int64  `form:"page" validate:"omitempty,min=1"`
	Limit int64  `form:"limit" validate:"omitempty,min=1,max=100"`
	Sort  string `form:"sort" validate:"omitempty,oneof=asc desc"`
}

type QueryParamGetMovieList struct {
	Keyword    string `form:"keyword" validate:"omitempty,max=200"`
	Category   string `form:"category" validate:"omitempty,max=200"`
	Pagination QueryParamPagination
}
