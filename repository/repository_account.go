package repository

import (
	"github.com/google/uuid"
	modelDatabase "github.com/parinyapt/prinflix_backend/model/database"
	modelRepository "github.com/parinyapt/prinflix_backend/model/repository"
	"github.com/pkg/errors"
)

func (receiver RepositoryReceiverArgument) CreateAccount(param modelRepository.ParamCreateAccount) (err error) {
	resultDB := receiver.databaseTX.Create(&modelDatabase.Account{
		UUID:          param.UUID,
		Name:          param.Name,
		Email:         param.Email,
		EmailVerified: param.EmailVerified,
		Password:      param.Password,
		Status:        param.Status,
		Role:          param.Role,
	})
	if resultDB.Error != nil {
		return errors.Wrap(resultDB.Error, "[Repository][CreateAccount()]->"+errorDatabaseQueryFailed)
	}

	return nil
}

func (receiver RepositoryReceiverArgument) FetchOneAccountByEmail(email string) (result modelRepository.ResultFetchOneAccount, err error) {
	resultDB := receiver.databaseTX.Where(&modelDatabase.Account{Email: email}).Limit(1).Find(&result.Data)
	if resultDB.Error != nil {
		return result, errors.Wrap(resultDB.Error, "[Repository][FetchOneAccountByEmail()]->"+errorDatabaseQueryFailed)
	}
	if resultDB.RowsAffected == 0 {
		return result, nil
	}

	result.IsFound = true

	return result, nil
}

func (receiver RepositoryReceiverArgument) FetchOneAccountByUUID(uuid uuid.UUID) (result modelRepository.ResultFetchOneAccount, err error) {
	resultDB := receiver.databaseTX.Where(&modelDatabase.Account{UUID: uuid}).Limit(1).Find(&result.Data)
	if resultDB.Error != nil {
		return result, errors.Wrap(resultDB.Error, "[Repository][FetchOneAccountByUUID()]->"+errorDatabaseQueryFailed)
	}
	if resultDB.RowsAffected == 0 {
		return result, nil
	}

	result.IsFound = true

	return result, nil
}

func (receiver RepositoryReceiverArgument) UpdateAccountByUUID(uuid uuid.UUID, param modelRepository.ParamUpdateAccount) (result modelRepository.ResultIsFoundOnly, err error) {
	resultDB := receiver.databaseTX.Model(&modelDatabase.Account{UUID: uuid}).Updates(&modelDatabase.Account{
		Name:          param.Name,
		Email:         param.Email,
		EmailVerified: param.EmailVerified,
		Password:      param.Password,
		Status:        param.Status,
		Image:         param.Image,
		Role:          param.Role,
	})
	if resultDB.Error != nil {
		return result, errors.Wrap(resultDB.Error, "[Repository][UpdateAccountByUUID()]->"+errorDatabaseQueryFailed)
	}
	if resultDB.RowsAffected == 0 {
		return result, nil
	}

	result.IsFound = true

	return result, nil
}

func (receiver RepositoryReceiverArgument) DeleteAccountByUUID(uuid uuid.UUID) (result modelRepository.ResultIsFoundOnly, err error) {
	resultDB := receiver.databaseTX.Where(&modelDatabase.Account{UUID: uuid}).Delete(&modelDatabase.Account{})
	if resultDB.Error != nil {
		return result, errors.Wrap(resultDB.Error, "[Repository][DeleteAccountByUUID()]->"+errorDatabaseQueryFailed)
	}
	if resultDB.RowsAffected == 0 {
		return result, nil
	}

	result.IsFound = true

	return result, nil
}
