package modelHandler

import (
	"github.com/google/uuid"
)

type RequestAddMovieComment struct {
	Comment string `json:"comment" validate:"required,min=1,max=1000,alpha_space"`
}

type UriParamDeleteMovieComment struct {
	MovieUUID   string `uri:"movie_uuid" validate:"required,uuid"`
	CommentUUID string `uri:"comment_uuid" validate:"required,uuid"`
}

type ResponseGetMovieComment struct {
	ResultPagination ResponsePagination         `json:"result_detail"`
	ResultData       []ResponseMovieCommentData `json:"result_data"`
}

type ResponseMovieCommentData struct {
	CommentUUID      uuid.UUID `json:"uuid"`
	AccountName      string    `json:"account_name"`
	CommentContent   string    `json:"comment"`
	CommentCreatedAt string    `json:"created_at"`
}
