package repository

import (
	modelDatabase "github.com/parinyapt/prinflix_backend/model/database"
	modelRepository "github.com/parinyapt/prinflix_backend/model/repository"
	"github.com/pkg/errors"
)

func (receiver RepositoryReceiverArgument) CreateAccountOAuth(param modelRepository.ParamCreateAccountOAuth) (err error) {
	resultDB := receiver.databaseTX.Create(&modelDatabase.AccountOAuth{
		AccountUUID: param.AccountUUID,
		Provider:    param.Provider,
		UserID:      param.UserID,
		UserEmail:   param.UserEmail,
		UserPicture: param.UserPicture,
	})
	if resultDB.Error != nil {
		return errors.Wrap(resultDB.Error, "[Repository][CreateAccountOAuth()]->"+errorDatabaseQueryFailed)
	}

	return nil
}

func (receiver RepositoryReceiverArgument) FetchOneAccountOAuthByProviderAndUserID(param modelRepository.ParamFetchOneAccountOAuthByProviderAndUserID) (result modelRepository.ResultFetchOneAccountOAuth, err error) {
	resultDB := receiver.databaseTX.Where(&modelDatabase.AccountOAuth{
		Provider:    param.Provider,
		UserID:      param.UserID,
	}).Limit(1).Find(&result.Data)
	if resultDB.Error != nil {
		return result, errors.Wrap(resultDB.Error, "[Repository][FetchOneAccountOAuthByProviderAndUserID()]->"+errorDatabaseQueryFailed)
	}
	if resultDB.RowsAffected == 0 {
		return result, nil
	}

	result.IsFound = true

	return result, nil
}

func (receiver RepositoryReceiverArgument) FetchOneAccountOAuthByProviderAndAccountUUID(param modelRepository.ParamFetchOneAccountOAuthByProviderAndAccountUUID) (result modelRepository.ResultFetchOneAccountOAuth, err error) {
	resultDB := receiver.databaseTX.Where(&modelDatabase.AccountOAuth{
		Provider:    param.Provider,
		AccountUUID: param.AccountUUID,
	}).Limit(1).Find(&result.Data)
	if resultDB.Error != nil {
		return result, errors.Wrap(resultDB.Error, "[Repository][FetchOneAccountOAuthByProviderAndAccountUUID()]->"+errorDatabaseQueryFailed)
	}
	if resultDB.RowsAffected == 0 {
		return result, nil
	}

	result.IsFound = true

	return result, nil
}

func (receiver RepositoryReceiverArgument) DeleteAccountOAuthByProviderAndAccountUUID(param modelRepository.ParamDeleteAccountOAuthByProviderAndAccountUUID) (result modelRepository.ResultIsFoundOnly, err error) {
	resultDB := receiver.databaseTX.Where(&modelDatabase.AccountOAuth{
		Provider:    param.Provider,
		AccountUUID: param.AccountUUID,
	}).Delete(&modelDatabase.AccountOAuth{})
	if resultDB.Error != nil {
		return result, errors.Wrap(resultDB.Error, "[Repository][DeleteAccountOAuthByProviderAndAccountUUID()]->"+errorDatabaseQueryFailed)
	}
	if resultDB.RowsAffected == 0 {
		return result, nil
	}

	result.IsFound = true

	return result, nil
}