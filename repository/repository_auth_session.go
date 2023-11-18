package repository

import (
	"github.com/google/uuid"
	modelDatabase "github.com/parinyapt/prinflix_backend/model/database"
	modelRepository "github.com/parinyapt/prinflix_backend/model/repository"
	"github.com/pkg/errors"
)

func (receiver RepositoryReceiverArgument) CreateAuthSession(param modelRepository.ParamCreateAuthSession) (err error) {
	resultDB := receiver.databaseTX.Create(&modelDatabase.AuthSession{
		UUID:        param.SessionUUID,
		AccountUUID: param.AccountUUID,
		ExpiredAt:   param.ExpiredAt,
	})
	if resultDB.Error != nil {
		return errors.Wrap(resultDB.Error, "[Repository][CreateAuthSession()]->"+errorDatabaseQueryFailed)
	}

	return nil
}

func (receiver RepositoryReceiverArgument) FetchOneAuthSessionByUUID(sessionUUID uuid.UUID) (result modelRepository.ResultFetchOneAuthSession, err error) {
	resultDB := receiver.databaseTX.Where(&modelDatabase.AuthSession{UUID: sessionUUID}).Limit(1).Find(&result.Data)
	if resultDB.Error != nil {
		return result, errors.Wrap(resultDB.Error, "[Repository][FetchOneAuthSessionByUUID()]->"+errorDatabaseQueryFailed)
	}
	if resultDB.RowsAffected == 0 {
		return result, nil
	}

	result.IsFound = true

	return result, nil
}

func (receiver RepositoryReceiverArgument) DeleteAuthSessionByAccountUUID(accountUUID uuid.UUID) (result modelRepository.ResultIsFoundOnly, err error) {
	resultDB := receiver.databaseTX.Where(&modelDatabase.AuthSession{AccountUUID: accountUUID}).Delete(&modelDatabase.AuthSession{})
	if resultDB.Error != nil {
		return result, errors.Wrap(resultDB.Error, "[Repository][DeleteAuthSessionByAccountUUID()]->"+errorDatabaseQueryFailed)
	}
	if resultDB.RowsAffected == 0 {
		return result, nil
	}

	result.IsFound = true

	return result, nil
}
