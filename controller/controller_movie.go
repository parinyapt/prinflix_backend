package controller

import (
	modelController "github.com/parinyapt/prinflix_backend/model/controller"
	modelRepository "github.com/parinyapt/prinflix_backend/model/repository"
	"github.com/parinyapt/prinflix_backend/repository"
	"github.com/parinyapt/prinflix_backend/storage"
	utilsUUID "github.com/parinyapt/prinflix_backend/utils/uuid"
	"github.com/pkg/errors"
)

func (receiver ControllerReceiverArgument) GetAllMovieCategory() (returnData modelController.ReturnGetAllMovieCategory, err error) {
	repoInstance := repository.NewRepository(receiver.databaseTX)

	repoData, repoErr := repoInstance.FetchManyMovieCategory()
	if repoErr != nil {
		return returnData, errors.Wrap(repoErr, "[Controller][GetAllMovieCategory()]->Fail to fetch many movie category")
	}
	if !repoData.IsFound {
		returnData.IsNotFound = true
		return returnData, nil
	}

	for _, data := range repoData.Data {
		returnData.Data = append(returnData.Data, modelController.ReturnGetAllMovieCategoryData{
			CategoryID:   data.ID,
			CategoryName: data.Name,
		})
	}

	return returnData, nil
}

func (receiver ControllerReceiverArgument) GetAllMovie(param modelController.ParamGetAllMovie) (returnData modelController.ReturnGetManyMovie, err error) {
	accountUUIDparse, err := utilsUUID.ParseUUIDfromString(param.AccountUUID)
	if err != nil {
		return returnData, errors.Wrap(err, "[Controller][GetAllMovie()]->Fail to parse account uuid")
	}

	repoInstance := repository.NewRepository(receiver.databaseTX)

	repoData, repoErr := repoInstance.FetchManyMovie(accountUUIDparse, modelRepository.ParamFetchManyMovie{
		CategoryID:  param.CategoryID,
		SearchQuery: param.SearchQuery,
		Pagination: modelRepository.ParamPagination{
			Page:        param.Pagination.Page,
			Limit:       param.Pagination.Limit,
			SortField:   param.Pagination.SortField,
			SortOrderBy: param.Pagination.SortOrderBy,
		},
	})
	if repoErr != nil {
		return returnData, errors.Wrap(repoErr, "[Controller][GetAllMovie()]->Fail to fetch many movie")
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
			IsFavorite:        data.IsFavorite,
		})
	}

	returnData.Pagination.TotalData = repoData.Pagination.TotalData
	returnData.Pagination.TotalPage = repoData.Pagination.TotalPage
	returnData.Pagination.Page = repoData.Pagination.Page
	returnData.Pagination.Limit = repoData.Pagination.Limit

	return returnData, nil
}
