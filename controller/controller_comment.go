package controller

import (
	modelController "github.com/parinyapt/prinflix_backend/model/controller"
	modelRepository "github.com/parinyapt/prinflix_backend/model/repository"
	"github.com/parinyapt/prinflix_backend/repository"
	utilsUUID "github.com/parinyapt/prinflix_backend/utils/uuid"
	"github.com/pkg/errors"
)

func (receiver ControllerReceiverArgument) CreateMovieComment(param modelController.ParamCreateMovieComment) (returnData modelController.ReturnCreateMovieComment, err error) {
	movieUUIDparse, err := utilsUUID.ParseUUIDfromString(param.MovieUUID)
	if err != nil {
		return returnData, errors.Wrap(err, "[Controller][CreateMovieComment()]->Fail to parse movie uuid")
	}

	accountUUIDparse, err := utilsUUID.ParseUUIDfromString(param.AccountUUID)
	if err != nil {
		return returnData, errors.Wrap(err, "[Controller][CreateMovieComment()]->Fail to parse account uuid")
	}

	repoInstance := repository.NewRepository(receiver.databaseTX)

	repoData, repoErr := repoInstance.FetchOneMovie(modelRepository.ParamFetchOneMovie{
		MovieUUID:   movieUUIDparse,
		AccountUUID: accountUUIDparse,
	})
	if repoErr != nil {
		return returnData, errors.Wrap(repoErr, "[Controller][CreateMovieComment()]->Fail to fetch one movie")
	}
	if !repoData.IsFound {
		returnData.IsNotFound = true
		return returnData, nil
	}

	commentUUID := utilsUUID.GenerateUUIDv4()
	repoErr = repoInstance.CreateMovieComment(modelRepository.ParamCreateMovieComment{
		CommentUUID: commentUUID,
		AccountUUID: accountUUIDparse,
		MovieUUID:   movieUUIDparse,
		Comment:     param.Comment,
	})
	if repoErr != nil {
		return returnData, errors.Wrap(repoErr, "[Controller][CreateMovieComment()]->Fail to create movie comment")
	}

	returnData.CommentUUID = commentUUID.String()

	return returnData, nil
}

func (receiver ControllerReceiverArgument) GetMovieComment(param modelController.ParamGetMovieComment) (returnData modelController.ReturnGetMovieComment, err error) {
	movieUUIDparse, err := utilsUUID.ParseUUIDfromString(param.MovieUUID)
	if err != nil {
		return returnData, errors.Wrap(err, "[Controller][GetMovieComment()]->Fail to parse movie uuid")
	}

	repoInstance := repository.NewRepository(receiver.databaseTX)

	repoData, repoErr := repoInstance.FetchManyMovieCommentByMovieUUID(movieUUIDparse, modelRepository.ParamFetchManyMovieComment{
		Pagination: modelRepository.ParamPagination{
			Page:        param.Pagination.Page,
			Limit:       param.Pagination.Limit,
			SortField:   param.Pagination.SortField,
			SortOrderBy: param.Pagination.SortOrderBy,
		},
	})
	if repoErr != nil {
		return returnData, errors.Wrap(repoErr, "[Controller][GetMovieComment()]->Fail to fetch many movie")
	}
	if !repoData.IsFound {
		returnData.IsNotFound = true
		return returnData, nil
	}

	for _, data := range repoData.Data {
		returnData.Data = append(returnData.Data, modelController.ReturnGetMovieCommentData{
			CommentUUID:      data.CommentUUID,
			AccountName:      data.AccountName,
			CommentContent:   data.CommentContent,
			CommentCreatedAt: data.CommentCreatedAt,
		})
	}

	returnData.Pagination.TotalData = repoData.Pagination.TotalData
	returnData.Pagination.TotalPage = repoData.Pagination.TotalPage
	returnData.Pagination.Page = repoData.Pagination.Page
	returnData.Pagination.Limit = repoData.Pagination.Limit

	return returnData, nil
}

func (receiver ControllerReceiverArgument) DeleteMovieComment(param modelController.ParamDeleteMovieComment) (returnData modelController.ReturnIsNotFoundOnly, err error) {
	commentUUIDparse, err := utilsUUID.ParseUUIDfromString(param.CommentUUID)
	if err != nil {
		return returnData, errors.Wrap(err, "[Controller][DeleteMovieComment()]->Fail to parse comment uuid")
	}

	accountUUIDparse, err := utilsUUID.ParseUUIDfromString(param.AccountUUID)
	if err != nil {
		return returnData, errors.Wrap(err, "[Controller][DeleteMovieComment()]->Fail to parse account uuid")
	}

	repoInstance := repository.NewRepository(receiver.databaseTX)

	repoData, repoErr := repoInstance.FetchOneMovieCommentByCommentUUID(commentUUIDparse)
	if repoErr != nil {
		return returnData, errors.Wrap(repoErr, "[Controller][DeleteMovieComment()]->Fail to fetch one movie comment by comment uuid")
	}
	if !repoData.IsFound {
		returnData.IsNotFound = true
		return returnData, nil
	}

	if repoData.Data.AccountUUID != accountUUIDparse {
		returnData.IsNotFound = true
		return returnData, nil
	}

	_, repoErr = repoInstance.DeleteMovieCommentByCommentUUID(commentUUIDparse)
	if repoErr != nil {
		return returnData, errors.Wrap(repoErr, "[Controller][DeleteMovieComment()]->Fail to delete movie comment")
	}

	return returnData, nil
}
