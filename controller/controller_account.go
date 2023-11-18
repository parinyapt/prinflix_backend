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

func (receiver ControllerReceiverArgument) GetAccountInfo(accountUUID string) (returnData modelController.ReturnGetAccountInfo, err error) {
	accountUUIDparse, err := utilsUUID.ParseUUIDfromString(accountUUID)
	if err != nil {
		return returnData, errors.Wrap(err, "[Controller][GetAccountInfo()]->Fail to parse account uuid")
	}
	
	repoInstance := repository.NewRepository(receiver.databaseTX)

	repoData, repoErr := repoInstance.FetchOneAccountByUUID(accountUUIDparse)
	if repoErr != nil {
		return returnData, errors.Wrap(repoErr, "[Controller][GetAccountInfo()]->Fail to fetch one account by uuid")
	}

	if !repoData.IsFound {
		returnData.IsNotFound = true
		return returnData, nil
	}

	returnData.Name = repoData.Data.Name
	returnData.Email = repoData.Data.Email
	returnData.EmailVerified = repoData.Data.EmailVerified
	returnData.Status = repoData.Data.Status
	returnData.Image = repoData.Data.Image
	returnData.Role = repoData.Data.Role

	return returnData, nil
}
