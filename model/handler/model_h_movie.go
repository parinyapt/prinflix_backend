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
	ResultPagination ResponsePagination  `json:"result_detail"`
	ResultData       []ResponseMovieData `json:"result_data"`
}

type ResponseMovieData struct {
	MovieUUID         uuid.UUID               `json:"uuid"`
	MovieThumbnail    string                  `json:"thumbnail"`
	MovieTitle        string                  `json:"title"`
	MovieDescription  string                  `json:"description"`
	MovieCategoryID   uint                    `json:"category_id"`
	MovieCategoryName string                  `json:"category_name"`
	IsFavorite        bool                    `json:"is_favorite"`
	Review            ResponseMovieDataReview `json:"review"`
}

type ResponseMovieDataReview struct {
	ReviewTotalCount     int64   `json:"total_count"`
	ReviewGoodCount      int64   `json:"good_count"`
	ReviewFairCount      int64   `json:"fair_count"`
	ReviewBadCount       int64   `json:"bad_count"`
	ReviewGoodPercentage float64 `json:"good_percentage"`
	ReviewFairPercentage float64 `json:"fair_percentage"`
	ReviewBadPercentage  float64 `json:"bad_percentage"`
}

type UriParamMovieUUIDonly struct {
	MovieUUID string `uri:"movie_uuid" validate:"required,uuid"`
}
