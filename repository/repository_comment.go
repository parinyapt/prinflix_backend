package repository

import (
	"fmt"
	"math"

	"github.com/google/uuid"
	PTGUdata "github.com/parinyapt/golang_utils/data/v1"
	modelDatabase "github.com/parinyapt/prinflix_backend/model/database"
	modelRepository "github.com/parinyapt/prinflix_backend/model/repository"
	utilsDatabase "github.com/parinyapt/prinflix_backend/utils/database"
	"github.com/pkg/errors"
)

func (receiver RepositoryReceiverArgument) CreateMovieComment(param modelRepository.ParamCreateMovieComment) (err error) {
	resultDB := receiver.databaseTX.Create(&modelDatabase.Comment{
		UUID:        param.CommentUUID,
		AccountUUID: param.AccountUUID,
		MovieUUID:   param.MovieUUID,
		Comment:     param.Comment,
	})
	if resultDB.Error != nil {
		return errors.Wrap(resultDB.Error, "[Repository][CreateMovieComment()]->"+errorDatabaseQueryFailed)
	}

	return nil
}

func (receiver RepositoryReceiverArgument) FetchOneMovieCommentByCommentUUID(commentUUID uuid.UUID) (result modelRepository.ResultFetchOneMovieComment, err error) {
	resultDB := receiver.databaseTX.Where(&modelDatabase.Comment{
		UUID: commentUUID,
	}).Limit(1).Find(&result.Data)
	if resultDB.Error != nil {
		return result, errors.Wrap(resultDB.Error, "[Repository][FetchOneMovieCommentByCommentUUID()]->"+errorDatabaseQueryFailed)
	}
	if resultDB.RowsAffected == 0 {
		return result, nil
	}

	result.IsFound = true

	return result, nil
}

const (
	FetchManyMovieCommentSortFieldCreatedAt = "comment_created_at"
)

func (receiver RepositoryReceiverArgument) FetchManyMovieCommentByMovieUUID(movieUUID uuid.UUID, param modelRepository.ParamFetchManyMovieComment) (result modelRepository.ResultFetchManyMovieComment, err error) {
	// Pagination calculation
	receiver.databaseTX.Model(&modelDatabase.Comment{}).Where(&modelDatabase.Comment{
		MovieUUID: movieUUID,
	}).Count(&result.Pagination.TotalData)
	if result.Pagination.TotalData == 0 {
		result.IsFound = false
		return result, nil
	}
	result.Pagination.TotalPage = int64(math.Ceil(float64(result.Pagination.TotalData) / float64(param.Pagination.Limit)))
	result.Pagination.Page = param.Pagination.Page
	if result.Pagination.Page > result.Pagination.TotalPage {
		result.IsFound = false
		return result, nil
	}
	if param.Pagination.Limit > 100 {
		param.Pagination.Limit = 100
	}
	result.Pagination.Limit = param.Pagination.Limit

	sqlCommand := `SELECT
		@table_comment.comment_uuid AS comment_uuid,
		@table_account.account_name AS account_name,
		@table_comment.comment_content AS comment_content,
		@table_comment.comment_created_at AS comment_created_at
	FROM
		@table_comment
	INNER JOIN(
		SELECT comment_uuid FROM @table_comment WHERE @table_comment.comment_movie_uuid = @movie_uuid ORDER BY @sort_field @sort_order_by LIMIT @limit OFFSET @offset
	) AS tmp USING(comment_uuid)
	INNER JOIN @table_account ON @table_comment.comment_account_uuid = @table_account.account_uuid
	WHERE @table_comment.comment_movie_uuid = @movie_uuid
	ORDER BY @sort_field @sort_order_by;`
	sqlCommand = PTGUdata.ReplaceString(sqlCommand, map[string]string{
		"@table_comment": utilsDatabase.GenerateTableName("comment"),
		"@table_account": utilsDatabase.GenerateTableName("account"),
		"@sort_field":    param.Pagination.SortField,
		"@sort_order_by": param.Pagination.SortOrderBy,
		"@limit":         fmt.Sprintf("%d", param.Pagination.Limit),
		"@offset":        fmt.Sprintf("%d", (param.Pagination.Page*param.Pagination.Limit)-param.Pagination.Limit),
	})
	resultDB := receiver.databaseTX.Raw(sqlCommand, map[string]interface{}{
		"movie_uuid": movieUUID,
	}).Scan(&result.Data)
	if resultDB.Error != nil {
		return result, errors.Wrap(resultDB.Error, "[Repository][FetchManyMovieCommentByMovieUUID()]->"+errorDatabaseQueryFailed)
	}
	if resultDB.RowsAffected == 0 {
		return result, nil
	}

	result.IsFound = true

	return result, nil
}

func (receiver RepositoryReceiverArgument) DeleteMovieCommentByCommentUUID(commentUUID uuid.UUID) (result modelRepository.ResultIsFoundOnly, err error) {
	resultDB := receiver.databaseTX.Where(&modelDatabase.Comment{
		UUID: commentUUID,
	}).Delete(&modelDatabase.Comment{})
	if resultDB.Error != nil {
		return result, errors.Wrap(resultDB.Error, "[Repository][DeleteMovieCommentByCommentUUID()]->"+errorDatabaseQueryFailed)
	}
	if resultDB.RowsAffected == 0 {
		return result, nil
	}

	result.IsFound = true

	return result, nil
}
