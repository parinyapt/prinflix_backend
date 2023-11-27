package repository

import (
	"github.com/google/uuid"
	modelDatabase "github.com/parinyapt/prinflix_backend/model/database"
	modelRepository "github.com/parinyapt/prinflix_backend/model/repository"
	"github.com/pkg/errors"
)

func (receiver RepositoryReceiverArgument) CreateWatchSession(param modelRepository.ParamCreateWatchSession) (err error) {
	resultDB := receiver.databaseTX.Create(&modelDatabase.WatchSession{
		UUID:        param.SessionUUID,
		AccountUUID: param.AccountUUID,
		MovieUUID:   param.MovieUUID,
		ExpiredAt:   param.ExpiredAt,
	})
	if resultDB.Error != nil {
		return errors.Wrap(resultDB.Error, "[Repository][CreateFavoriteMovie()]->"+errorDatabaseQueryFailed)
	}

	return nil
}

func (receiver RepositoryReceiverArgument) FetchOneWatchSessionBySessionUUID(sessionUUID uuid.UUID) (result modelRepository.ResultFetchOneWatchSession, err error) {
	resultDB := receiver.databaseTX.Where(&modelDatabase.WatchSession{UUID: sessionUUID}).First(&result.Data)
	if resultDB.Error != nil {
		return result, errors.Wrap(resultDB.Error, "[Repository][FetchOneWatchSessionBySessionUUID()]->"+errorDatabaseQueryFailed)
	}
	if resultDB.RowsAffected == 0 {
		return result, nil
	}

	result.IsFound = true

	return result, nil
}

func (receiver RepositoryReceiverArgument) DeleteWatchSessionByAccountUUID(accountUUID uuid.UUID) (result modelRepository.ResultIsFoundOnly, err error) {
	resultDB := receiver.databaseTX.Where(&modelDatabase.WatchSession{AccountUUID: accountUUID}).Delete(&modelDatabase.WatchSession{})
	if resultDB.Error != nil {
		return result, errors.Wrap(resultDB.Error, "[Repository][DeleteWatchSessionByUUID()]->"+errorDatabaseQueryFailed)
	}
	if resultDB.RowsAffected == 0 {
		return result, nil
	}

	result.IsFound = true

	return result, nil
}
