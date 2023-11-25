package controller

import (
	modelController "github.com/parinyapt/prinflix_backend/model/controller"
	modelRepository "github.com/parinyapt/prinflix_backend/model/repository"
	"github.com/parinyapt/prinflix_backend/repository"
	utilsUUID "github.com/parinyapt/prinflix_backend/utils/uuid"

	PTGUpassword "github.com/parinyapt/golang_utils/password/v1"

	"github.com/pkg/errors"
)

func (receiver ControllerReceiverArgument) CreateAccount(param modelController.ParamCreateAccount) (returnData modelController.ReturnCreateAccount, err error) {
	repoInstance := repository.NewRepository(receiver.databaseTX)

	repoData, repoErr := repoInstance.FetchOneAccountByEmail(param.Email)
	if repoErr != nil {
		return returnData, errors.Wrap(repoErr, "[Controller][CreateAccount()]->Fail to fetch one account by email")
	}

	if repoData.IsFound {
		returnData.IsExist = true
		return returnData, nil
	}

	passwordHash, err := PTGUpassword.HashPassword(param.Password, 14)
	if err != nil {
		return returnData, errors.Wrap(err, "[Controller][CreateAccount()]->Fail to hash password")
	}

	accountUUID := utilsUUID.GenerateUUIDv4()
	repoErr = repoInstance.CreateAccount(modelRepository.ParamCreateAccount{
		UUID:          accountUUID,
		Name:          param.Name,
		Email:         param.Email,
		EmailVerified: false,
		Password:      passwordHash,
		Status:        "active",
		Role:          "user",
	})
	if repoErr != nil {
		return returnData, errors.Wrap(repoErr, "[Controller][CreateAccount()]->Fail to create account")
	}

	returnData.UUID = accountUUID

	return returnData, nil
}

func (receiver ControllerReceiverArgument) GetAccountInfo(param modelController.ParamGetAccountInfo) (returnData modelController.ReturnGetAccountInfo, err error) {

	if param.AccountUUID == "" && param.Email == "" {
		return returnData, errors.New("[Controller][GetAccountInfo()]->Both account uuid and email are empty")
	}

	repoInstance := repository.NewRepository(receiver.databaseTX)
	var repoData modelRepository.ResultFetchOneAccount
	var repoErr error

	if param.AccountUUID != "" {
		accountUUIDparse, err := utilsUUID.ParseUUIDfromString(param.AccountUUID)
		if err != nil {
			return returnData, errors.Wrap(err, "[Controller][GetAccountInfo()]->Fail to parse account uuid")
		}
		repoData, repoErr = repoInstance.FetchOneAccountByUUID(accountUUIDparse)
	} else {
		repoData, repoErr = repoInstance.FetchOneAccountByEmail(param.Email)
	}

	if repoErr != nil {
		return returnData, errors.Wrap(repoErr, "[Controller][GetAccountInfo()]->Fail to fetch one account by uuid")
	}

	if !repoData.IsFound {
		returnData.IsNotFound = true
		return returnData, nil
	}

	returnData.AccountUUID = repoData.Data.UUID
	returnData.Name = repoData.Data.Name
	returnData.Email = repoData.Data.Email
	returnData.EmailVerified = repoData.Data.EmailVerified
	returnData.Status = repoData.Data.Status
	returnData.Image = repoData.Data.Image
	returnData.Role = repoData.Data.Role

	return returnData, nil
}

func (receiver ControllerReceiverArgument) UpdateAccount(accountUUID string, param modelController.ParamUpdateAccount) (returnData modelController.ReturnIsNotFoundOnly, err error) {
	accountUUIDparse, err := utilsUUID.ParseUUIDfromString(accountUUID)
	if err != nil {
		return returnData, errors.Wrap(err, "[Controller][UpdateAccount()]->Fail to parse account uuid")
	}

	repoInstance := repository.NewRepository(receiver.databaseTX)

	if param.Password != "" {
		passwordHash, err := PTGUpassword.HashPassword(param.Password, 14)
		if err != nil {
			return returnData, errors.Wrap(err, "[Controller][UpdateAccount()]->Fail to hash password")
		}
		param.Password = passwordHash
	}

	repoData, repoErr := repoInstance.UpdateAccountByUUID(accountUUIDparse, modelRepository.ParamUpdateAccount{
		Name:          param.Name,
		EmailVerified: param.EmailVerified,
		Password:      param.Password,
		Image:         param.Image,
	})
	if repoErr != nil {
		return returnData, errors.Wrap(repoErr, "[Controller][UpdateAccount()]->Fail to update account")
	}

	if !repoData.IsFound {
		returnData.IsNotFound = true
		return returnData, nil
	}

	return returnData, nil
}
