package modelRepository

import (
	"github.com/google/uuid"
)

type ParamFetchOneMovie struct {
	AccountUUID uuid.UUID
	MovieUUID   uuid.UUID
}

type ResultFetchOneMovie struct {
	IsFound bool
	Data    DBResultFetchMovie
}

type ParamFetchManyMovie struct {
	SearchQuery string
	CategoryID  uint
	Pagination  ParamPagination
}

type ResultFetchManyMovie struct {
	IsFound    bool
	Data       []DBResultFetchMovie
	Pagination ResultPagination
}

type DBResultFetchMovie struct {
	MovieUUID         uuid.UUID
	MovieTitle        string
	MovieDescription  string
	MovieCategoryID   uint
	MovieCategoryName string
	IsFavorite        bool
}
