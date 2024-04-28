package controller

import (
	"time"

	modelController "github.com/parinyapt/prinflix_backend/model/controller"
	modelRepository "github.com/parinyapt/prinflix_backend/model/repository"
	"github.com/parinyapt/prinflix_backend/repository"
	utilsUUID "github.com/parinyapt/prinflix_backend/utils/uuid"

	"github.com/pkg/errors"
)

const (
	OauthStateExpiredIn = time.Minute * 15
)

func (receiver ControllerReceiverArgument) CreateOauthState(provider string) (returnData modelController.ReturnCreateOauthState, err error) {
	repoInstance := repository.NewRepository(receiver.databaseTX)

	stateUUID := utilsUUID.GenerateUUIDv4()
	repoErr := repoInstance.CreateOauthState(modelRepository.ParamCreateOauthState{
		UUID:     stateUUID,
		Provider: provider,
	})
	if repoErr != nil {
		return returnData, errors.Wrap(repoErr, "[Controller][CreateOauthState()]->Fail to create oauth state")
	}

	returnData.StateUUID = stateUUID

	return returnData, nil
}

func (receiver ControllerReceiverArgument) CheckOauthState(param modelController.ParamOauthState) (returnData modelController.ReturnCheckOauthState, err error) {
	stateUUIDparse, err := utilsUUID.ParseUUIDfromString(param.StateUUID)
	if err != nil {
		return returnData, errors.Wrap(err, "[Controller][CheckOauthState()]->Fail to parse state uuid")
	}

	repoInstance := repository.NewRepository(receiver.databaseTX)

	repoData, repoErr := repoInstance.FetchOneOauthState(stateUUIDparse)
	if repoErr != nil {
		return returnData, errors.Wrap(repoErr, "[Controller][CheckOauthState()]->Fail to fetch one oauth state by uuid")
	}

	if !repoData.IsFound {
		returnData.IsNotFound = true
		return returnData, nil
	}

	if repoData.Data.Provider != param.Provider {
		returnData.IsNotFound = true
		return returnData, nil
	}

	if time.Now().After(repoData.Data.CreatedAt.Add(OauthStateExpiredIn)) {
		returnData.IsExpired = true
		return returnData, nil
	}

	returnData.StateUUID = repoData.Data.UUID

	return returnData, nil
}

func (receiver ControllerReceiverArgument) DeleteOauthState(param modelController.ParamOauthState) (err error) {
	stateUUIDparse, err := utilsUUID.ParseUUIDfromString(param.StateUUID)
	if err != nil {
		return errors.Wrap(err, "[Controller][DeleteOauthState()]->Fail to parse state uuid")
	}

	repoInstance := repository.NewRepository(receiver.databaseTX)

	_, repoErr := repoInstance.DeleteOauthStateByUUID(stateUUIDparse)

	if repoErr != nil {
		return errors.Wrap(repoErr, "[Controller][DeleteOauthState()]->Fail to delete oauth state by uuid")
	}

	return nil
}
