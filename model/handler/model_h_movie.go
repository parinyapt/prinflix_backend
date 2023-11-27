package modelHandler

import "github.com/google/uuid"

type ResponseGetMovieCategoryList struct {
	CategoryID   uint   `json:"id"`
	CategoryName string `json:"name"`
}

type QueryParamGetMovieList struct {
	Keyword    string               `form:"keyword" validate:"omitempty,max=200"`
	CategoryID uint                 `form:"category" validate:"omitempty,number"`
	Pagination QueryParamPagination `form:"pagination"`
}

type ResponseGetMovieList struct {
	ResultPagination ResponsePagination         `json:"result_detail"`
	ResultData       []ResponseMovieData `json:"result_data"`
}

type ResponseMovieData struct {
	MovieUUID         uuid.UUID `json:"uuid"`
	MovieThumbnail    string    `json:"thumbnail"`
	MovieTitle        string    `json:"title"`
	MovieDescription  string    `json:"description"`
	MovieCategoryID   uint      `json:"category_id"`
	MovieCategoryName string    `json:"category_name"`
	IsFavorite        bool      `json:"is_favorite"`
}

type UriParamGetMovieDetail struct {
	MovieUUID string `uri:"movie_uuid" binding:"required,uuid"`
}