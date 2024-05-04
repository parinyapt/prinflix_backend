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
		returnData.Data = append(returnData.Data, modelController.ReturnGetMovieCategoryData{
			CategoryID:   data.ID,
			CategoryName: data.Name,
		})
	}

	return returnData, nil
}

func (receiver ControllerReceiverArgument) GetMovieCategoryDetail(categoryID uint) (returnData modelController.ReturnGetMovieCategoryDetail, err error) {
	repoInstance := repository.NewRepository(receiver.databaseTX)

	repoData, repoErr := repoInstance.FetchOneMovieCategoryByID(categoryID)
	if repoErr != nil {
		return returnData, errors.Wrap(repoErr, "[Controller][GetMovieCategoryDetail()]->Fail to fetch one movie category by id")
	}
	if !repoData.IsFound {
		returnData.IsNotFound = true
		return returnData, nil
	}

	returnData.Data.CategoryID = repoData.Data.ID
	returnData.Data.CategoryName = repoData.Data.Name

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

		var reviewGoodPercentage float64
		var reviewFairPercentage float64
		var reviewBadPercentage float64
		if data.ReviewTotalCount == 0 {
			reviewGoodPercentage = 0
			reviewFairPercentage = 0
			reviewBadPercentage = 0
		} else {
			reviewGoodPercentage = float64((float64(data.ReviewGoodCount) / float64(data.ReviewTotalCount)) * 100)
			reviewFairPercentage = float64((float64(data.ReviewFairCount) / float64(data.ReviewTotalCount)) * 100)
			reviewBadPercentage = float64((float64(data.ReviewBadCount) / float64(data.ReviewTotalCount)) * 100)
		}

		returnData.Data = append(returnData.Data, modelController.ReturnGetManyMovieData{
			MovieUUID:            data.MovieUUID,
			MovieThumbnail:       thumbnailURL,
			MovieTitle:           data.MovieTitle,
			MovieDescription:     data.MovieDescription,
			MovieCategoryID:      data.MovieCategoryID,
			MovieCategoryName:    data.MovieCategoryName,
			IsFavorite:           data.IsFavorite,
			ReviewTotalCount:     data.ReviewTotalCount,
			ReviewGoodCount:      data.ReviewGoodCount,
			ReviewFairCount:      data.ReviewFairCount,
			ReviewBadCount:       data.ReviewBadCount,
			ReviewGoodPercentage: reviewGoodPercentage,
			ReviewFairPercentage: reviewFairPercentage,
			ReviewBadPercentage:  reviewBadPercentage,
		})
	}

	returnData.Pagination.TotalData = repoData.Pagination.TotalData
	returnData.Pagination.TotalPage = repoData.Pagination.TotalPage
	returnData.Pagination.Page = repoData.Pagination.Page
	returnData.Pagination.Limit = repoData.Pagination.Limit

	return returnData, nil
}

func (receiver ControllerReceiverArgument) GetMovieDetail(param modelController.ParamAccountUUIDandMovieUUID) (returnData modelController.ReturnGetMovieDetail, err error) {
	movieUUIDparse, err := utilsUUID.ParseUUIDfromString(param.MovieUUID)
	if err != nil {
		return returnData, errors.Wrap(err, "[Controller][GetMovieDetail()]->Fail to parse movie uuid")
	}

	accountUUIDparse, err := utilsUUID.ParseUUIDfromString(param.AccountUUID)
	if err != nil {
		return returnData, errors.Wrap(err, "[Controller][GetMovieDetail()]->Fail to parse account uuid")
	}

	repoInstance := repository.NewRepository(receiver.databaseTX)

	repoData, repoErr := repoInstance.FetchOneMovie(modelRepository.ParamFetchOneMovie{
		MovieUUID:   movieUUIDparse,
		AccountUUID: accountUUIDparse,
	})
	if repoErr != nil {
		return returnData, errors.Wrap(repoErr, "[Controller][GetMovieDetail()]->Fail to fetch one movie")
	}
	if !repoData.IsFound {
		returnData.IsNotFound = true
		return returnData, nil
	}

	returnData.MovieUUID = repoData.Data.MovieUUID
	returnData.MovieTitle = repoData.Data.MovieTitle
	returnData.MovieDescription = repoData.Data.MovieDescription
	returnData.MovieCategoryID = repoData.Data.MovieCategoryID
	returnData.MovieCategoryName = repoData.Data.MovieCategoryName
	returnData.MovieThumbnail, err = storage.GenerateRoutePath(1, storage.MovieThumbnailRoutePath, map[string]string{
		"movie_uuid": repoData.Data.MovieUUID.String(),
	})
	if err != nil {
		returnData.MovieThumbnail = ""
	}
	returnData.ReviewTotalCount = repoData.Data.ReviewTotalCount
	returnData.ReviewGoodCount = repoData.Data.ReviewGoodCount
	returnData.ReviewFairCount = repoData.Data.ReviewFairCount
	returnData.ReviewBadCount = repoData.Data.ReviewBadCount
	if repoData.Data.ReviewTotalCount == 0 {
		returnData.ReviewGoodPercentage = 0
		returnData.ReviewFairPercentage = 0
		returnData.ReviewBadPercentage = 0
	} else {
		returnData.ReviewGoodPercentage = float64((float64(repoData.Data.ReviewGoodCount) / float64(repoData.Data.ReviewTotalCount)) * 100)
		returnData.ReviewFairPercentage = float64((float64(repoData.Data.ReviewFairCount) / float64(repoData.Data.ReviewTotalCount)) * 100)
		returnData.ReviewBadPercentage = float64((float64(repoData.Data.ReviewBadCount) / float64(repoData.Data.ReviewTotalCount)) * 100)
	}

	return returnData, nil
}

func (receiver ControllerReceiverArgument) GetAllRecommendMovie(accountUUID string) (returnData modelController.ReturnGetManyMovie, err error) {
	accountUUIDparse, err := utilsUUID.ParseUUIDfromString(accountUUID)
	if err != nil {
		return returnData, errors.Wrap(err, "[Controller][GetAllRecommendMovie()]->Fail to parse account uuid")
	}

	repoInstance := repository.NewRepository(receiver.databaseTX)

	repoData, repoErr := repoInstance.FetchManyMovie(accountUUIDparse, modelRepository.ParamFetchManyMovie{
		Pagination: modelRepository.ParamPagination{
			Page:        1,
			Limit:       5,
			SortField:   "RAND()",
			SortOrderBy: "",
		},
	})
	if repoErr != nil {
		return returnData, errors.Wrap(repoErr, "[Controller][GetAllRecommendMovie()]->Fail to fetch many movie")
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

func (receiver ControllerReceiverArgument) GetRecommendMovieByMostViewCategory(accountUUID string) (returnData modelController.ReturnGetManyMovie, err error) {
	accountUUIDparse, err := utilsUUID.ParseUUIDfromString(accountUUID)
	if err != nil {
		return returnData, errors.Wrap(err, "[Controller][GetRecommendMovieByMostViewCategory()]->Fail to parse account uuid")
	}

	repoInstance := repository.NewRepository(receiver.databaseTX)

	repoDataMostViewCategory, repoErr := repoInstance.FetchOneAccountMostViewCategory(accountUUIDparse)
	if repoErr != nil {
		return returnData, errors.Wrap(repoErr, "[Controller][GetRecommendMovieByMostViewCategory()]->Fail to fetch one account most view category")
	}
	if !repoDataMostViewCategory.IsFound {
		returnData.IsNotFound = true
		return returnData, nil
	}

	repoData, repoErr := repoInstance.FetchManyMovie(accountUUIDparse, modelRepository.ParamFetchManyMovie{
		Pagination: modelRepository.ParamPagination{
			Page:        1,
			Limit:       5,
			SortField:   "RAND()",
			SortOrderBy: "",
		},
		CategoryID: repoDataMostViewCategory.Data.CategoryId,
	})
	if repoErr != nil {
		return returnData, errors.Wrap(repoErr, "[Controller][GetRecommendMovieByMostViewCategory()]->Fail to fetch many movie")
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
