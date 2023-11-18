package repository

import (
	modelDatabase "github.com/parinyapt/prinflix_backend/model/database"
	modelRepository "github.com/parinyapt/prinflix_backend/model/repository"
	"github.com/pkg/errors"
)

func (receiver RepositoryReceiverArgument) FetchOneMovieCategoryByID(id uint) (result modelRepository.ResultFetchOneMovieCategory, err error) {
	resultDB := receiver.databaseTX.Where(&modelDatabase.MovieCategory{ID: id}).Limit(1).Find(&result.Data)
	if resultDB.Error != nil {
		return result, errors.Wrap(resultDB.Error, "[Repository][FetchOneMovieCategoryByID()]->"+errorDatabaseQueryFailed)
	}
	if resultDB.RowsAffected == 0 {
		return result, nil
	}

	result.IsFound = true

	return result, nil
}

func (receiver RepositoryReceiverArgument) FetchManyMovieCategory() (result modelRepository.ResultFetchManyMovieCategory, err error) {
	resultDB := receiver.databaseTX.Find(&result.Data)
	if resultDB.Error != nil {
		return result, errors.Wrap(resultDB.Error, "[Repository][FetchManyMovieCategory()]->"+errorDatabaseQueryFailed)
	}
	if resultDB.RowsAffected == 0 {
		return result, nil
	}

	result.IsFound = true

	return result, nil
}
