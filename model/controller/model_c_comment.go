package modelController

import (
	"time"

	"github.com/google/uuid"
)

type ParamCreateMovieComment struct {
	AccountUUID string
	MovieUUID   string
	Comment     string
}

type ReturnCreateMovieComment struct {
	IsNotFound  bool
	CommentUUID string
}

type ParamDeleteMovieComment struct {
	AccountUUID string
	CommentUUID string
}

type ParamGetMovieComment struct {
	MovieUUID  string
	Pagination ParamPagination
}

type ReturnGetMovieComment struct {
	IsNotFound bool
	Data       []ReturnGetMovieCommentData
	Pagination ReturnPagination
}

type ReturnGetMovieCommentData struct {
	CommentUUID      uuid.UUID
	AccountName      string
	CommentContent   string
	CommentCreatedAt time.Time
}
