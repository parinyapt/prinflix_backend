package controller

import (
	modelController "github.com/parinyapt/prinflix_backend/model/controller"
	modelRepository "github.com/parinyapt/prinflix_backend/model/repository"
	"github.com/parinyapt/prinflix_backend/repository"
	utilsUUID "github.com/parinyapt/prinflix_backend/utils/uuid"

	"github.com/pkg/errors"
)

func (receiver ControllerReceiverArgument) CreateAccountOAuth(param modelController.ParamCreateAccountOAuth) (err error) {
	accountUUIDparse, err := utilsUUID.ParseUUIDfromString(param.AccountUUID)
	if err != nil {
		return errors.Wrap(err, "[Controller][CreateAccountOAuth()]->Fail to parse account uuid")
	}

	repoInstance := repository.NewRepository(receiver.databaseTX)

	repoErr := repoInstance.CreateAccountOAuth(modelRepository.ParamCreateAccountOAuth{
		AccountUUID: accountUUIDparse,
		Provider:    param.Provider,
		UserID:      param.UserID,
		UserName:    param.UserName,
		UserEmail:   param.UserEmail,
		UserPicture: param.UserPicture,
	})
	if repoErr != nil {
		return errors.Wrap(repoErr, "[Controller][CreateAccountOAuth()]->Fail to create account oauth")
	}

	return nil
}

func (receiver ControllerReceiverArgument) CheckAccountOAuth(provider string, param modelController.ParamCheckAccountOAuth) (returnData modelController.ReturnCheckAccountOAuth, err error) {
	repoInstance := repository.NewRepository(receiver.databaseTX)

	var repoData modelRepository.ResultFetchOneAccountOAuth
	var repoErr error

	if param.AccountUUID != "" {
		accountUUIDparse, err := utilsUUID.ParseUUIDfromString(param.AccountUUID)
		if err != nil {
			return returnData, errors.Wrap(err, "[Controller][CheckAccountOAuth()]->Fail to parse account uuid")
		}

		repoData, repoErr = repoInstance.FetchOneAccountOAuthByProviderAndAccountUUID(modelRepository.ParamFetchOneAccountOAuthByProviderAndAccountUUID{
			Provider:    provider,
			AccountUUID: accountUUIDparse,
		})
	}else {
		repoData, repoErr = repoInstance.FetchOneAccountOAuthByProviderAndUserID(modelRepository.ParamFetchOneAccountOAuthByProviderAndUserID{
			Provider: provider,
			UserID:   param.UserID,
		})
	}

	if repoErr != nil {
		return returnData, errors.Wrap(repoErr, "[Controller][CheckAccountOAuth()]->Fail to fetch one account oauth by provider and user id or account uuid")
	}

	if !repoData.IsFound {
		returnData.IsNotFound = true
		return returnData, nil
	}

	returnData.AccountUUID = repoData.Data.AccountUUID
	returnData.Name = repoData.Data.UserName
	returnData.Email = repoData.Data.UserEmail
	returnData.Picture = repoData.Data.UserPicture

	return returnData, nil
}

func (receiver ControllerReceiverArgument) DeleteAccountOAuth(param modelController.ParamDeleteAccountOAuth) (returnData modelController.ReturnIsNotFoundOnly, err error) {
	accountUUIDparse, err := utilsUUID.ParseUUIDfromString(param.AccountUUID)
	if err != nil {
		return returnData, errors.Wrap(err, "[Controller][DeleteAccountOAuth()]->Fail to parse account uuid")
	}

	repoInstance := repository.NewRepository(receiver.databaseTX)

	repoDeleteData, repoErr := repoInstance.DeleteAccountOAuthByProviderAndAccountUUID(modelRepository.ParamDeleteAccountOAuthByProviderAndAccountUUID{
		AccountUUID: accountUUIDparse,
		Provider: param.Provider,
	})
	if repoErr != nil {
		return returnData, errors.Wrap(repoErr, "[Controller][DeleteAccountOAuth()]->Fail to delete account oauth by provider and account uuid")
	}
	if !repoDeleteData.IsFound {
		returnData.IsNotFound = true
		return returnData, nil
	}

	return returnData, nil
}