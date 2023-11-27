package modelController

import "github.com/google/uuid"

type ReturnGetAllMovieCategory struct {
	IsNotFound bool
	Data       []ReturnGetAllMovieCategoryData
}

type ReturnGetAllMovieCategoryData struct {
	CategoryID   uint   `json:"id"`
	CategoryName string `json:"name"`
}

type ParamGetAllMovie struct {
	AccountUUID string
	SearchQuery string
	CategoryID  uint
	Pagination  ParamPagination
}

type ReturnGetManyMovie struct {
	IsNotFound bool
	Data       []ReturnGetManyMovieData
	Pagination ReturnPagination
}

type ReturnGetManyMovieData struct {
	MovieUUID         uuid.UUID
	MovieThumbnail    string
	MovieTitle        string
	MovieDescription  string
	MovieCategoryID   uint
	MovieCategoryName string
	IsFavorite        bool
}

type ReturnGetMovieDetail struct {
	IsNotFound bool

	MovieUUID         uuid.UUID
	MovieThumbnail    string
	MovieTitle        string
	MovieDescription  string
	MovieCategoryID   uint
	MovieCategoryName string
	IsFavorite        bool
}

type ParamGetMovieDetail struct {
	AccountUUID string
	MovieUUID   string
}
