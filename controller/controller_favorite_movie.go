package controller

import (
	modelController "github.com/parinyapt/prinflix_backend/model/controller"
	modelRepository "github.com/parinyapt/prinflix_backend/model/repository"
	"github.com/parinyapt/prinflix_backend/repository"
	"github.com/parinyapt/prinflix_backend/storage"
	utilsUUID "github.com/parinyapt/prinflix_backend/utils/uuid"
	"github.com/pkg/errors"
)

func (receiver ControllerReceiverArgument) CreateFavoriteMovie(param modelController.ParamAccountUUIDandMovieUUID) (err error) {
	accountUUIDparse, err := utilsUUID.ParseUUIDfromString(param.AccountUUID)
	if err != nil {
		return errors.Wrap(err, "[Controller][CreateFavoriteMovie()]->Fail to parse account uuid")
	}

	movieUUIDparse, err := utilsUUID.ParseUUIDfromString(param.MovieUUID)
	if err != nil {
		return errors.Wrap(err, "[Controller][CreateFavoriteMovie()]->Fail to parse movie uuid")
	}

	repoInstance := repository.NewRepository(receiver.databaseTX)

	repoErr := repoInstance.CreateFavoriteMovie(modelRepository.ParamFavoriteMovie{
		AccountUUID: accountUUIDparse,
		MovieUUID:   movieUUIDparse,
	})
	if repoErr != nil {
		return errors.Wrap(repoErr, "[Controller][CreateAccountOAuth()]->Fail to create account oauth")
	}

	return nil
}

func (receiver ControllerReceiverArgument) CheckFavoriteMovie(param modelController.ParamAccountUUIDandMovieUUID) (returnData modelController.ReturnIsNotFoundOnly, err error) {
	accountUUIDparse, err := utilsUUID.ParseUUIDfromString(param.AccountUUID)
	if err != nil {
		return returnData, errors.Wrap(err, "[Controller][CheckFavoriteMovie()]->Fail to parse account uuid")
	}

	movieUUIDparse, err := utilsUUID.ParseUUIDfromString(param.MovieUUID)
	if err != nil {
		return returnData, errors.Wrap(err, "[Controller][CheckFavoriteMovie()]->Fail to parse movie uuid")
	}

	repoInstance := repository.NewRepository(receiver.databaseTX)

	repoData, repoErr := repoInstance.FetchOneFavoriteMovieByAccountUUIDAndMovieUUID(modelRepository.ParamFavoriteMovie{
		AccountUUID: accountUUIDparse,
		MovieUUID:   movieUUIDparse,
	})
	if repoErr != nil {
		return returnData, errors.Wrap(repoErr, "[Controller][CheckFavoriteMovie()]->Fail to fetch one favorite movie by account uuid and movie uuid")
	}
	if !repoData.IsFound {
		returnData.IsNotFound = true
		return returnData, nil
	}

	return returnData, nil
}

func (receiver ControllerReceiverArgument) DeleteFavoriteMovie(param modelController.ParamAccountUUIDandMovieUUID) (returnData modelController.ReturnIsNotFoundOnly, err error) {
	accountUUIDparse, err := utilsUUID.ParseUUIDfromString(param.AccountUUID)
	if err != nil {
		return returnData, errors.Wrap(err, "[Controller][DeleteFavoriteMovie()]->Fail to parse account uuid")
	}

	movieUUIDparse, err := utilsUUID.ParseUUIDfromString(param.MovieUUID)
	if err != nil {
		return returnData, errors.Wrap(err, "[Controller][DeleteFavoriteMovie()]->Fail to parse movie uuid")
	}

	repoInstance := repository.NewRepository(receiver.databaseTX)

	repoData, repoErr := repoInstance.DeleteFavoriteMovieByAccountUUIDAndMovieUUID(modelRepository.ParamFavoriteMovie{
		AccountUUID: accountUUIDparse,
		MovieUUID:   movieUUIDparse,
	})
	if repoErr != nil {
		return returnData, errors.Wrap(repoErr, "[Controller][DeleteFavoriteMovie()]->Fail to delete favorite movie by account uuid and movie uuid")
	}
	if !repoData.IsFound {
		returnData.IsNotFound = true
		return returnData, nil
	}

	return returnData, nil
}

func (receiver ControllerReceiverArgument) GetAllFavoriteMovie(accountUUID string) (returnData modelController.ReturnGetAllFavoriteMovie, err error) {
	accountUUIDparse, err := utilsUUID.ParseUUIDfromString(accountUUID)
	if err != nil {
		return returnData, errors.Wrap(err, "[Controller][GetAllFavoriteMovie()]->Fail to parse account uuid")
	}

	repoInstance := repository.NewRepository(receiver.databaseTX)

	repoData, repoErr := repoInstance.FetchManyFavoriteMovieByAccountUUID(accountUUIDparse)
	if repoErr != nil {
		return returnData, errors.Wrap(repoErr, "[Controller][GetAllFavoriteMovie()]->Fail to fetch many favorite movie by account uuid")
	}
	if !repoData.IsFound {
		returnData.IsNotFound = true
		return returnData, nil
	}

	for _, data := range repoData.Data {
		thumbnailURL, err := storage.GenerateRoutePath(1, storage.MovieThumbnailRoutePath, map[string]string{
			"movie_uuid": data.MovieUUID.String(),
		})
		if err != nil {
			thumbnailURL = ""
		}

		returnData.Data = append(returnData.Data, modelController.ReturnGetManyMovieData{
			MovieUUID:         data.MovieUUID,
			MovieThumbnail:    thumbnailURL,
			MovieTitle:        data.MovieTitle,
			MovieDescription:  data.MovieDescription,
			MovieCategoryID:   data.MovieCategoryID,
			MovieCategoryName: data.MovieCategoryName,
			IsFavorite:        true,
		})
	}

	return returnData, nil
}