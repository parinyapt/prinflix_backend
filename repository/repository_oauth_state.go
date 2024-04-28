package repository

import (
	"github.com/google/uuid"
	modelDatabase "github.com/parinyapt/prinflix_backend/model/database"
	modelRepository "github.com/parinyapt/prinflix_backend/model/repository"
	"github.com/pkg/errors"
)

func (receiver RepositoryReceiverArgument) CreateOauthState(param modelRepository.ParamCreateOauthState) (err error) {
	resultDB := receiver.databaseTX.Create(&modelDatabase.OauthState{
		UUID:     param.UUID,
		Provider: param.Provider,
	})
	if resultDB.Error != nil {
		return errors.Wrap(resultDB.Error, "[Repository][CreateOauthState()]->"+errorDatabaseQueryFailed)
	}

	return nil
}

func (receiver RepositoryReceiverArgument) FetchOneOauthState(uuid uuid.UUID) (result modelRepository.ResultFetchOneOauthState, err error) {
	resultDB := receiver.databaseTX.Where(&modelDatabase.OauthState{UUID: uuid}).Limit(1).Find(&result.Data)
	if resultDB.Error != nil {
		return result, errors.Wrap(resultDB.Error, "[Repository][FetchOneOauthState()]->"+errorDatabaseQueryFailed)
	}
	if resultDB.RowsAffected == 0 {
		return result, nil
	}

	result.IsFound = true

	return result, nil
}

func (receiver RepositoryReceiverArgument) DeleteOauthStateByUUID(uuid uuid.UUID) (result modelRepository.ResultIsFoundOnly, err error) {
	resultDB := receiver.databaseTX.Where(&modelDatabase.OauthState{UUID: uuid}).Delete(&modelDatabase.OauthState{})
	if resultDB.Error != nil {
		return result, errors.Wrap(resultDB.Error, "[Repository][DeleteOauthStateByUUID()]->"+errorDatabaseQueryFailed)
	}
	if resultDB.RowsAffected == 0 {
		return result, nil
	}

	result.IsFound = true

	return result, nil
}
