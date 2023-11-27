package modelRepository

import (
	"github.com/google/uuid"
	modelDatabase "github.com/parinyapt/prinflix_backend/model/database"
)

type ResultFetchOneMovie struct {
	IsFound bool
	Data    *modelDatabase.Movie
}

type ParamFetchManyMovie struct {
	SearchQuery string
	CategoryID  uint
	Pagination  ParamPagination
}

type ResultFetchManyMovie struct {
	IsFound    bool
	Data       []DBResultFetchManyMovie
	Pagination ResultPagination
}

type DBResultFetchManyMovie struct {
	MovieUUID         uuid.UUID
	MovieTitle        string
	MovieDescription  string
	MovieCategoryID   uint
	MovieCategoryName string
	IsFavorite        bool
}
