package controller

import (
	modelController "github.com/parinyapt/prinflix_backend/model/controller"
	"github.com/parinyapt/prinflix_backend/repository"

	PTGUpassword "github.com/parinyapt/golang_utils/password/v1"

	"github.com/pkg/errors"
)

func (receiver ControllerReceiverArgument) CheckLogin(param modelController.ParamCheckLogin) (returnData modelController.ReturnCheckLogin, err error) {
	repoInstance := repository.NewRepository(receiver.databaseTX)

	repoData, repoErr := repoInstance.FetchOneAccountByEmail(param.Email)
	if repoErr != nil {
		return returnData, errors.Wrap(repoErr, "[Controller][CreateAccount()]->Fail to fetch one account by email")
	}

	if repoData.IsFound {
		returnData.IsNotFound = true
		return returnData, nil
	}

	isPasswordMatch := PTGUpassword.VerifyHashPassword(param.Password, repoData.Data.Password)
	if !isPasswordMatch {
		returnData.IsPasswordNotMatch = true
		return returnData, nil
	}

	returnData.UUID = repoData.Data.UUID

	return returnData, nil
}
