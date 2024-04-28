package modelController

import "github.com/google/uuid"

type ReturnGetMovieCategoryDetail struct {
	IsNotFound bool
	Data       ReturnGetMovieCategoryData
}

type ReturnGetAllMovieCategory struct {
	IsNotFound bool
	Data       []ReturnGetMovieCategoryData
}

type ReturnGetMovieCategoryData struct {
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
	MovieUUID            uuid.UUID
	MovieThumbnail       string
	MovieTitle           string
	MovieDescription     string
	MovieCategoryID      uint
	MovieCategoryName    string
	IsFavorite           bool
	ReviewTotalCount     int64
	ReviewGoodCount      int64
	ReviewFairCount      int64
	ReviewBadCount       int64
	ReviewGoodPercentage float64
	ReviewFairPercentage float64
	ReviewBadPercentage  float64
}

type ReturnGetMovieDetail struct {
	IsNotFound bool

	MovieUUID            uuid.UUID
	MovieThumbnail       string
	MovieTitle           string
	MovieDescription     string
	MovieCategoryID      uint
	MovieCategoryName    string
	IsFavorite           bool
	ReviewTotalCount     int64
	ReviewGoodCount      int64
	ReviewFairCount      int64
	ReviewBadCount       int64
	ReviewGoodPercentage float64
	ReviewFairPercentage float64
	ReviewBadPercentage  float64
}

type ParamAccountUUIDandMovieUUID struct {
	AccountUUID string
	MovieUUID   string
}
