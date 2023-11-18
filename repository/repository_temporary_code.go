package repository

import (
	"github.com/google/uuid"
	modelDatabase "github.com/parinyapt/prinflix_backend/model/database"
	modelRepository "github.com/parinyapt/prinflix_backend/model/repository"
	"github.com/pkg/errors"
)

func (receiver RepositoryReceiverArgument) CreateTemporaryCode(param modelRepository.ParamCreateTemporaryCode) (err error) {
	resultDB := receiver.databaseTX.Create(&modelDatabase.TemporaryCode{
		UUID:        param.UUID,
		AccountUUID: param.AccountUUID,
		Type:        param.Type,
	})
	if resultDB.Error != nil {
		return errors.Wrap(resultDB.Error, "[Repository][CreateTemporaryCode()]->"+errorDatabaseQueryFailed)
	}

	return nil
}

func (receiver RepositoryReceiverArgument) FetchOneTemporaryCodeByUUID(uuid uuid.UUID) (result modelRepository.ResultFetchOneTemporaryCode, err error) {
	resultDB := receiver.databaseTX.Where(&modelDatabase.TemporaryCode{UUID: uuid}).Limit(1).Find(&result.Data)
	if resultDB.Error != nil {
		return result, errors.Wrap(resultDB.Error, "[Repository][FetchOneTemporaryCodeByUUID()]->"+errorDatabaseQueryFailed)
	}
	if resultDB.RowsAffected == 0 {
		return result, nil
	}

	result.IsFound = true

	return result, nil
}

func (receiver RepositoryReceiverArgument) DeleteTemporaryCodeByUUID(uuid uuid.UUID) (result modelRepository.ResultIsFoundOnly, err error) {
	resultDB := receiver.databaseTX.Where(&modelDatabase.TemporaryCode{UUID: uuid}).Delete(&modelDatabase.TemporaryCode{})
	if resultDB.Error != nil {
		return result, errors.Wrap(resultDB.Error, "[Repository][DeleteTemporaryCodeByUUID()]->"+errorDatabaseQueryFailed)
	}
	if resultDB.RowsAffected == 0 {
		return result, nil
	}

	result.IsFound = true

	return result, nil
}

func (receiver RepositoryReceiverArgument) DeleteTemporaryCodeByAccountUUIDAndType(param modelRepository.ParamDeleteTemporaryCodeByAccountUUIDAndType) (result modelRepository.ResultIsFoundOnly, err error) {
	resultDB := receiver.databaseTX.Where(&modelDatabase.TemporaryCode{
		AccountUUID: param.AccountUUID,
		Type:        param.Type,
	}).Delete(&modelDatabase.TemporaryCode{})
	if resultDB.Error != nil {
		return result, errors.Wrap(resultDB.Error, "[Repository][DeleteTemporaryCodeByAccountUUIDAndType()]->"+errorDatabaseQueryFailed)
	}
	if resultDB.RowsAffected == 0 {
		return result, nil
	}

	result.IsFound = true

	return result, nil
}
