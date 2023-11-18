package controller

import (
	"time"

	modelController "github.com/parinyapt/prinflix_backend/model/controller"
	modelRepository "github.com/parinyapt/prinflix_backend/model/repository"
	"github.com/parinyapt/prinflix_backend/repository"
	utilsUUID "github.com/parinyapt/prinflix_backend/utils/uuid"

	"github.com/pkg/errors"
)

func (receiver ControllerReceiverArgument) CreateAuthSession(accountUUID string) (returnData modelController.ReturnCreateAuthSession, err error) {
	accountUUIDparse, err := utilsUUID.ParseUUIDfromString(accountUUID)
	if err != nil {
		return returnData, errors.Wrap(err, "[Controller][CreateAuthSession()]->Fail to parse account uuid")
	}

	repoInstance := repository.NewRepository(receiver.databaseTX)

	sessionUUID := utilsUUID.GenerateUUIDv4()
	sessionExpireIn := time.Hour * 24 * 15

	repoErr := repoInstance.CreateAuthSession(modelRepository.ParamCreateAuthSession{
		SessionUUID: sessionUUID,
		AccountUUID: accountUUIDparse,
		ExpiredAt:   time.Now().Add(sessionExpireIn),
	})
	if repoErr != nil {
		return returnData, errors.Wrap(repoErr, "[Controller][CreateAuthSession()]->Fail to create auth session")
	}

	returnData.SessionUUID = sessionUUID
	returnData.ExpiredAt = time.Now().Add(sessionExpireIn)
	returnData.ExtiredIn = sessionExpireIn

	return returnData, nil
}

func (receiver ControllerReceiverArgument) CheckAuthSession(sessionUUID string) (returnData modelController.ReturnCheckAuthSession, err error) {
	sessionUUIDparse, err := utilsUUID.ParseUUIDfromString(sessionUUID)
	if err != nil {
		return returnData, errors.Wrap(err, "[Controller][CheckAuthSession()]->Fail to parse session uuid")
	}

	repoInstance := repository.NewRepository(receiver.databaseTX)

	repoData, repoErr := repoInstance.FetchOneAuthSessionByUUID(sessionUUIDparse)
	if repoErr != nil {
		return returnData, errors.Wrap(repoErr, "[Controller][CheckAuthSession()]->Fail to fetch one auth session by uuid")
	}

	if !repoData.IsFound {
		returnData.IsNotFound = true
		return returnData, nil
	}

	if repoData.Data.ExpiredAt.Before(time.Now()) {
		returnData.IsExpired = true
		return returnData, nil
	}

	returnData.AccountUUID = repoData.Data.AccountUUID
	returnData.SessionUUID = repoData.Data.UUID

	return returnData, nil
}

func (receiver ControllerReceiverArgument) DeleteAuthSessionByAccountUUID(accountUUID string) (err error) {
	accountUUIDparse, err := utilsUUID.ParseUUIDfromString(accountUUID)
	if err != nil {
		return errors.Wrap(err, "[Controller][DeleteAuthSessionByAccountUUID()]->Fail to parse account uuid")
	}

	repoInstance := repository.NewRepository(receiver.databaseTX)

	_, repoErr := repoInstance.DeleteAuthSessionByAccountUUID(accountUUIDparse)
	if repoErr != nil {
		return errors.Wrap(repoErr, "[Controller][DeleteAuthSessionByAccountUUID()]->Fail to delete auth session")
	}

	return nil
}
