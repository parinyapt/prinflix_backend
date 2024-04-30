package modelRepository

import (
	"time"

	"github.com/google/uuid"
	modelDatabase "github.com/parinyapt/prinflix_backend/model/database"
)

type ParamCreateMovieComment struct {
	CommentUUID uuid.UUID
	AccountUUID uuid.UUID
	MovieUUID   uuid.UUID
	Comment     string
}

type ResultFetchOneMovieComment struct {
	IsFound bool
	Data    *modelDatabase.Comment
}

type ParamFetchManyMovieComment struct {
	Pagination ParamPagination
}

type ResultFetchManyMovieComment struct {
	IsFound    bool
	Data       []DBResultFetchMovieComment
	Pagination ResultPagination
}

type DBResultFetchMovieComment struct {
	CommentUUID      uuid.UUID
	AccountName      string
	CommentContent   string
	CommentCreatedAt time.Time
}
