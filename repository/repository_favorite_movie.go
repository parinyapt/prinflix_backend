package repository

import (
	"github.com/google/uuid"
	modelDatabase "github.com/parinyapt/prinflix_backend/model/database"
	modelRepository "github.com/parinyapt/prinflix_backend/model/repository"
	"github.com/pkg/errors"
)

func (receiver RepositoryReceiverArgument) CreateFavoriteMovie(param modelRepository.ParamFavoriteMovie) (err error) {
	resultDB := receiver.databaseTX.Create(&modelDatabase.FavoriteMovie{
		AccountUUID: param.AccountUUID,
		MovieUUID:   param.MovieUUID,
	})
	if resultDB.Error != nil {
		return errors.Wrap(resultDB.Error, "[Repository][CreateFavoriteMovie()]->"+errorDatabaseQueryFailed)
	}

	return nil
}

func (receiver RepositoryReceiverArgument) FetchOneFavoriteMovieByAccountUUIDAndMovieUUID(param modelRepository.ParamFavoriteMovie) (result modelRepository.ResultFetchOneFavoriteMovie, err error) {
	resultDB := receiver.databaseTX.Where(&modelDatabase.FavoriteMovie{
		AccountUUID: param.AccountUUID,
		MovieUUID:   param.MovieUUID,
	}).Limit(1).Find(&result.Data)
	if resultDB.Error != nil {
		return result, errors.Wrap(resultDB.Error, "[Repository][FetchOneFavoriteMovieByAccountUUIDAndMovieUUID()]->"+errorDatabaseQueryFailed)
	}
	if resultDB.RowsAffected == 0 {
		return result, nil
	}

	result.IsFound = true

	return result, nil
}

func (receiver RepositoryReceiverArgument) FetchManyFavoriteMovieByAccountUUID(accountUUID uuid.UUID) (result modelRepository.ResultFetchManyFavoriteMovie, err error) {
	resultDB := receiver.databaseTX.Where(&modelDatabase.FavoriteMovie{AccountUUID: accountUUID}).Find(&result.Data)
	if resultDB.Error != nil {
		return result, errors.Wrap(resultDB.Error, "[Repository][FetchManyFavoriteMovieByAccountUUID()]->"+errorDatabaseQueryFailed)
	}
	if resultDB.RowsAffected == 0 {
		return result, nil
	}

	return result, nil
}

func (receiver RepositoryReceiverArgument) DeleteFavoriteMovieByAccountUUIDAndMovieUUID(param modelRepository.ParamFavoriteMovie) (result modelRepository.ResultIsFoundOnly, err error) {
	resultDB := receiver.databaseTX.Where(&modelDatabase.FavoriteMovie{
		AccountUUID: param.AccountUUID,
		MovieUUID:   param.MovieUUID,
	}).Delete(&modelDatabase.FavoriteMovie{})
	if resultDB.Error != nil {
		return result, errors.Wrap(resultDB.Error, "[Repository][DeleteFavoriteMovieByAccountUUIDAndMovieUUID()]->"+errorDatabaseQueryFailed)
	}
	if resultDB.RowsAffected == 0 {
		return result, nil
	}

	result.IsFound = true

	return result, nil
}