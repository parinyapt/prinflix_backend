package repository

import (
	"github.com/google/uuid"
	modelDatabase "github.com/parinyapt/prinflix_backend/model/database"
	modelRepository "github.com/parinyapt/prinflix_backend/model/repository"
	"github.com/pkg/errors"
)

func (receiver RepositoryReceiverArgument) FetchOneMovie(movieUUID uuid.UUID) (result modelRepository.ResultFetchOneMovie, err error) {
	resultDB := receiver.databaseTX.Where(&modelDatabase.Movie{UUID: movieUUID}).Limit(1).Find(&result.Data)
	if resultDB.Error != nil {
		return result, errors.Wrap(resultDB.Error, "[Repository][FetchOneMovie()]->"+errorDatabaseQueryFailed)
	}
	if resultDB.RowsAffected == 0 {
		return result, nil
	}

	result.IsFound = true

	return result, nil
}

// func (receiver RepositoryReceiverArgument) FetchManyMovie() (result modelRepository.ResultFetchManyMovie, err error) {}