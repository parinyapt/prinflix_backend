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

func (receiver RepositoryReceiverArgument) FetchOneMovie(param modelRepository.ParamFetchOneMovie) (result modelRepository.ResultFetchOneMovie, err error) {
	sqlCommand := `SELECT
		movie_uuid,
		movie_title,
		movie_description,
		movie_category_id,
		movie_category AS movie_category_name,
		CASE
			WHEN favorite_movie_created_at is not NULL THEN true
			ELSE false
		END AS is_favorite
	FROM
		@table_movie
	INNER JOIN @table_movie_category ON movie_movie_category_id = movie_category_id
	LEFT JOIN(
		SELECT
			favorite_movie_movie_uuid
		FROM
			@table_favorite_movie
		WHERE
			favorite_movie_account_uuid = @account_uuid
	) tt ON movie_uuid = tt.favorite_movie_movie_uuid
	LEFT JOIN @table_favorite_movie ON tt.favorite_movie_movie_uuid = @table_favorite_movie.favorite_movie_movie_uuid
	WHERE
	(
		@table_favorite_movie.favorite_movie_account_uuid = @account_uuid
		OR 
		@table_favorite_movie.favorite_movie_account_uuid IS NULL
	) AND movie_uuid = @movie_uuid;`
	sqlCommand = PTGUdata.ReplaceString(sqlCommand, map[string]string{
		"@table_movie":          utilsDatabase.GenerateTableName("movie"),
		"@table_movie_category": utilsDatabase.GenerateTableName("movie_category"),
		"@table_favorite_movie": utilsDatabase.GenerateTableName("favorite_movie"),
	})
	resultDB := receiver.databaseTX.Raw(sqlCommand, map[string]interface{}{
		"account_uuid": param.AccountUUID,
		"movie_uuid":   param.MovieUUID,
	}).Scan(&result.Data)
	if resultDB.Error != nil {
		return result, errors.Wrap(resultDB.Error, "[Repository][FetchManyMovie()]->"+errorDatabaseQueryFailed)
	}
	if resultDB.RowsAffected == 0 {
		return result, nil
	}

	result.IsFound = true

	return result, nil
}

const (
	FetchManyMovieSortFieldMovieTitle = "movie_title"
)

func (receiver RepositoryReceiverArgument) FetchManyMovie(accountUUID uuid.UUID, param modelRepository.ParamFetchManyMovie) (result modelRepository.ResultFetchManyMovie, err error) {
	var whereCondition string
	whereCondition += "movie_title LIKE @search_keyword"
	if param.CategoryID > 0 {
		if len(whereCondition) > 0 {
			whereCondition += " AND "
		}
		whereCondition += "movie_movie_category_id = @category_id"
	}

	// Pagination calculation
	receiver.databaseTX.Model(&modelDatabase.Movie{}).Where(whereCondition, map[string]interface{}{
		"category_id":    param.CategoryID,
		"search_keyword": "%" + param.SearchQuery + "%",
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
		movie_uuid,
		movie_title,
		movie_description,
		movie_category_id,
		movie_category AS movie_category_name,
		CASE
			WHEN favorite_movie_created_at is not NULL THEN true
			ELSE false
		END AS is_favorite
	FROM
		@table_movie
	INNER JOIN(
		SELECT movie_uuid FROM @table_movie WHERE ` + whereCondition + ` ORDER BY @sort_field @sort_order_by LIMIT @limit OFFSET @offset
	) AS tmp USING(movie_uuid)
	INNER JOIN @table_movie_category ON movie_movie_category_id = movie_category_id
	LEFT JOIN(
		SELECT
			favorite_movie_movie_uuid
		FROM
			@table_favorite_movie
		WHERE
			favorite_movie_account_uuid = @account_uuid
	) tt ON movie_uuid = tt.favorite_movie_movie_uuid
	LEFT JOIN @table_favorite_movie ON tt.favorite_movie_movie_uuid = @table_favorite_movie.favorite_movie_movie_uuid
	WHERE
	(
		@table_favorite_movie.favorite_movie_account_uuid = @account_uuid
		OR 
		@table_favorite_movie.favorite_movie_account_uuid IS NULL
	)
	ORDER BY @sort_field @sort_order_by;`
	sqlCommand = PTGUdata.ReplaceString(sqlCommand, map[string]string{
		"@table_movie":          utilsDatabase.GenerateTableName("movie"),
		"@table_movie_category": utilsDatabase.GenerateTableName("movie_category"),
		"@table_favorite_movie": utilsDatabase.GenerateTableName("favorite_movie"),
		"@sort_field":           param.Pagination.SortField,
		"@sort_order_by":        param.Pagination.SortOrderBy,
		"@limit":                fmt.Sprintf("%d", param.Pagination.Limit),
		"@offset":               fmt.Sprintf("%d", (param.Pagination.Page*param.Pagination.Limit)-param.Pagination.Limit),
	})
	resultDB := receiver.databaseTX.Raw(sqlCommand, map[string]interface{}{
		"category_id":    param.CategoryID,
		"search_keyword": "%" + param.SearchQuery + "%",
		"account_uuid":   accountUUID,
	}).Scan(&result.Data)
	if resultDB.Error != nil {
		return result, errors.Wrap(resultDB.Error, "[Repository][FetchManyMovie()]->"+errorDatabaseQueryFailed)
	}
	if resultDB.RowsAffected == 0 {
		return result, nil
	}

	result.IsFound = true

	return result, nil
}