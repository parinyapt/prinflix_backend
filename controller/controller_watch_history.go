package controller

import (
	modelController "github.com/parinyapt/prinflix_backend/model/controller"
	modelRepository "github.com/parinyapt/prinflix_backend/model/repository"
	"github.com/parinyapt/prinflix_backend/repository"
	utilsUUID "github.com/parinyapt/prinflix_backend/utils/uuid"
	"github.com/pkg/errors"
)

func (receiver ControllerReceiverArgument) CreateWatchHistory(param modelController.ParamAccountUUIDandMovieUUID) (err error) {
	accountUUIDparse, err := utilsUUID.ParseUUIDfromString(param.AccountUUID)
	if err != nil {
		return errors.Wrap(err, "[Controller][CreateFavoriteMovie()]->Fail to parse account uuid")
	}

	movieUUIDparse, err := utilsUUID.ParseUUIDfromString(param.MovieUUID)
	if err != nil {
		return errors.Wrap(err, "[Controller][CreateFavoriteMovie()]->Fail to parse movie uuid")
	}

	repoInstance := repository.NewRepository(receiver.databaseTX)

	repoErr := repoInstance.CreateWatchHistory(modelRepository.ParamCreateWatchHistory{
		WatchHistoryUUID: utilsUUID.GenerateUUIDv4(),
		AccountUUID:      accountUUIDparse,
		MovieUUID:        movieUUIDparse,
	})
	if repoErr != nil {
		return errors.Wrap(repoErr, "[Controller][CreateWatchSession()]->Fail to create watch session")
	}

	return nil
}

func (receiver ControllerReceiverArgument) ClearWatchHistoryStart(param modelController.ParamUpdateWatchHistory) (returnData modelController.ReturnIsNotFoundOnly, err error) {
	accountUUIDparse, err := utilsUUID.ParseUUIDfromString(param.AccountUUID)
	if err != nil {
		return returnData, errors.Wrap(err, "[Controller][ClearWatchHistoryStart()]->Fail to parse account uuid")
	}

	movieUUIDparse, err := utilsUUID.ParseUUIDfromString(param.MovieUUID)
	if err != nil {
		return returnData, errors.Wrap(err, "[Controller][ClearWatchHistoryStart()]->Fail to parse movie uuid")
	}

	repoInstance := repository.NewRepository(receiver.databaseTX)

	repoData, repoErr := repoInstance.UpdateWatchHistoryIsEnd(modelRepository.ParamUpdateWatchHistory{
		AccountUUID: accountUUIDparse,
		MovieUUID:   movieUUIDparse,
		IsEnd:       false,
	})
	if repoErr != nil {
		return returnData, errors.Wrap(repoErr, "[Controller][ClearWatchHistoryStart()]->Fail to update watch history by account uuid and movie uuid")
	}

	if !repoData.IsFound {
		returnData.IsNotFound = true
		return returnData, nil
	}

	return returnData, nil
}

func (receiver ControllerReceiverArgument) ClearWatchHistoryPause(param modelController.ParamUpdateWatchHistory) (returnData modelController.ReturnIsNotFoundOnly, err error) {
	accountUUIDparse, err := utilsUUID.ParseUUIDfromString(param.AccountUUID)
	if err != nil {
		return returnData, errors.Wrap(err, "[Controller][ClearWatchHistoryPause()]->Fail to parse account uuid")
	}

	movieUUIDparse, err := utilsUUID.ParseUUIDfromString(param.MovieUUID)
	if err != nil {
		return returnData, errors.Wrap(err, "[Controller][ClearWatchHistoryPause()]->Fail to parse movie uuid")
	}

	repoInstance := repository.NewRepository(receiver.databaseTX)

	repoData, repoErr := repoInstance.UpdateWatchHistoryLatestTimeStamp(modelRepository.ParamUpdateWatchHistory{
		AccountUUID:     accountUUIDparse,
		MovieUUID:       movieUUIDparse,
		LatestTimeStamp: param.TimeStamp,
	})
	if repoErr != nil {
		return returnData, errors.Wrap(repoErr, "[Controller][ClearWatchHistoryPause()]->Fail to update watch history by account uuid and movie uuid")
	}

	if !repoData.IsFound {
		returnData.IsNotFound = true
		return returnData, nil
	}

	return returnData, nil
}

func (receiver ControllerReceiverArgument) ClearWatchHistoryEnd(param modelController.ParamUpdateWatchHistory) (returnData modelController.ReturnIsNotFoundOnly, err error) {
	accountUUIDparse, err := utilsUUID.ParseUUIDfromString(param.AccountUUID)
	if err != nil {
		return returnData, errors.Wrap(err, "[Controller][ClearWatchHistoryEnd()]->Fail to parse account uuid")
	}

	movieUUIDparse, err := utilsUUID.ParseUUIDfromString(param.MovieUUID)
	if err != nil {
		return returnData, errors.Wrap(err, "[Controller][ClearWatchHistoryEnd()]->Fail to parse movie uuid")
	}

	repoInstance := repository.NewRepository(receiver.databaseTX)

	repoDataIsEnd, repoErr := repoInstance.UpdateWatchHistoryIsEnd(modelRepository.ParamUpdateWatchHistory{
		AccountUUID: accountUUIDparse,
		MovieUUID:   movieUUIDparse,
		IsEnd:       true,
	})
	if repoErr != nil {
		return returnData, errors.Wrap(repoErr, "[Controller][ClearWatchHistoryEnd()]->Fail to update watch history by account uuid and movie uuid")
	}
	if !repoDataIsEnd.IsFound {
		returnData.IsNotFound = true
		return returnData, nil
	}

	repoDataLatestTimeStamp, repoErr := repoInstance.UpdateWatchHistoryLatestTimeStamp(modelRepository.ParamUpdateWatchHistory{
		AccountUUID:     accountUUIDparse,
		MovieUUID:       movieUUIDparse,
		LatestTimeStamp: 0,
	})
	if repoErr != nil {
		return returnData, errors.Wrap(repoErr, "[Controller][ClearWatchHistoryEnd()]->Fail to update watch history by account uuid and movie uuid")
	}
	if !repoDataLatestTimeStamp.IsFound {
		returnData.IsNotFound = true
		return returnData, nil
	}

	return returnData, nil
}

func (receiver ControllerReceiverArgument) CheckWatchHistory(param modelController.ParamAccountUUIDandMovieUUID) (returnData modelController.ReturnCheckWatchHistory, err error) {
	accountUUIDparse, err := utilsUUID.ParseUUIDfromString(param.AccountUUID)
	if err != nil {
		return returnData, errors.Wrap(err, "[Controller][CheckWatchHistory()]->Fail to parse account uuid")
	}

	movieUUIDparse, err := utilsUUID.ParseUUIDfromString(param.MovieUUID)
	if err != nil {
		return returnData, errors.Wrap(err, "[Controller][CheckWatchHistory()]->Fail to parse movie uuid")
	}

	repoInstance := repository.NewRepository(receiver.databaseTX)

	repoData, repoErr := repoInstance.FetchOneWatchHistory(modelRepository.ParamFetchOneWatchHistory{
		AccountUUID: accountUUIDparse,
		MovieUUID:   movieUUIDparse,
	})
	if repoErr != nil {
		return returnData, errors.Wrap(repoErr, "[Controller][CheckWatchHistory()]->Fail to fetch one watch session by session uuid")
	}

	if !repoData.IsFound {
		returnData.IsNotFound = true
		return returnData, nil
	}

	returnData.LatestTimeStamp = repoData.Data.LatestTimeStamp
	returnData.IsEnd = repoData.Data.IsEnd

	return returnData, nil
}
