package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	PTGUvalidator "github.com/parinyapt/golang_utils/validator/v1"
	"github.com/parinyapt/prinflix_backend/controller"
	"github.com/parinyapt/prinflix_backend/database"
	"github.com/parinyapt/prinflix_backend/logger"
	modelController "github.com/parinyapt/prinflix_backend/model/controller"
	modelHandler "github.com/parinyapt/prinflix_backend/model/handler"
	modelUtils "github.com/parinyapt/prinflix_backend/model/utils"
	"github.com/parinyapt/prinflix_backend/repository"
	utilsResponse "github.com/parinyapt/prinflix_backend/utils/response"
)

func GetMovieCommentHandler(c *gin.Context) {
	var uriParam modelHandler.UriParamMovieUUIDonly

	if err := c.ShouldBindUri(&uriParam); err != nil {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusBadRequest,
		})
		return
	}

	isValidatePass, _, validatorError := PTGUvalidator.Validate(uriParam)
	if validatorError != nil {
		logger.Error("[Handler][GetMovieCommentHandler()]->Error Validate Data", logger.Field("error", validatorError.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}
	if !isValidatePass {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusBadRequest,
			Error:        "Invalid Parameter",
		})
		return
	}

	var queryParam modelHandler.QueryParamPagination

	if err := c.ShouldBind(&queryParam); err != nil {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusBadRequest,
		})
		return
	}

	isValidatePass, _, validatorError = PTGUvalidator.Validate(queryParam)
	if validatorError != nil {
		logger.Error("[Handler][GetMovieCommentHandler()]->Error Validate Data", logger.Field("error", validatorError.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}
	if !isValidatePass {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusBadRequest,
			Error:        "Invalid Parameter",
		})
		return
	}

	if queryParam.Page == 0 {
		queryParam.Page = 1
	}

	if queryParam.Limit == 0 {
		queryParam.Limit = 10
	}

	queryParam.SortField = repository.FetchManyMovieCommentSortFieldCreatedAt

	if len(queryParam.SortOrderBy) == 0 {
		queryParam.SortOrderBy = repository.SortOrderByDesc
	}

	controllerInstance := controller.NewController(database.DB)

	getMovieComment, err := controllerInstance.GetMovieComment(modelController.ParamGetMovieComment{
		MovieUUID: uriParam.MovieUUID,
		Pagination: modelController.ParamPagination{
			Page:        queryParam.Page,
			Limit:       queryParam.Limit,
			SortField:   queryParam.SortField,
			SortOrderBy: queryParam.SortOrderBy,
		},
	})
	if err != nil {
		logger.Error("[Handler][GetMovieCommentHandler()]->Error GetAllMovie()", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}
	if getMovieComment.IsNotFound {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusNotFound,
			Error:        "Movie Comment Not Found",
		})
		return
	}

	var response modelHandler.ResponseGetMovieComment
	response.ResultPagination.TotalData = getMovieComment.Pagination.TotalData
	response.ResultPagination.TotalPage = getMovieComment.Pagination.TotalPage
	response.ResultPagination.CurrentPage = getMovieComment.Pagination.Page
	response.ResultPagination.CurrentLimit = getMovieComment.Pagination.Limit

	for _, comment := range getMovieComment.Data {
		response.ResultData = append(response.ResultData, modelHandler.ResponseMovieCommentData{
			CommentUUID:      comment.CommentUUID,
			AccountName:      comment.AccountName,
			CommentContent:   comment.CommentContent,
			CommentCreatedAt: comment.CommentCreatedAt.Format(time.RFC3339),
		})
	}

	utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
		ResponseCode: http.StatusOK,
		Data:         response,
	})
}

func AddMovieCommentHandler(c *gin.Context) {
	var uriParam modelHandler.UriParamMovieUUIDonly

	if err := c.ShouldBindUri(&uriParam); err != nil {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusBadRequest,
		})
		return
	}

	isValidatePass, _, validatorError := PTGUvalidator.Validate(uriParam)
	if validatorError != nil {
		logger.Error("[Handler][AddMovieCommentHandler()]->Error Validate Data", logger.Field("error", validatorError.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}
	if !isValidatePass {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusBadRequest,
			Error:        "Invalid Parameter",
		})
		return
	}

	var request modelHandler.RequestAddMovieComment
	if err := c.ShouldBindJSON(&request); err != nil {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusBadRequest,
		})
		return
	}

	isValidatePass, errorFieldList, validatorError := PTGUvalidator.Validate(request)
	if validatorError != nil {
		logger.Error("[Handler][AddMovieCommentHandler()]->Error Validate Data", logger.Field("error", validatorError.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}
	if !isValidatePass {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusBadRequest,
			Error:        errorFieldList,
		})
		return
	}

	databaseTx := database.DB.Begin()
	controllerInstance := controller.NewController(databaseTx)
	defer databaseTx.Rollback()

	addComment, err := controllerInstance.CreateMovieComment(modelController.ParamCreateMovieComment{
		AccountUUID: c.GetString("ACCOUNT_UUID"),
		MovieUUID:   uriParam.MovieUUID,
		Comment:     request.Comment,
	})
	if err != nil {
		logger.Error("[Handler][AddMovieCommentHandler()]->Error CreateMovieComment()", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}
	if addComment.IsNotFound {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusNotFound,
			Error:        "Movie Not Found",
		})
		return
	}

	databaseTx.Commit()

	utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
		ResponseCode: http.StatusOK,
		Data:         "Comment Success",
	})
}

func DeleteMovieCommentHandler(c *gin.Context) {
	var uriParam modelHandler.UriParamDeleteMovieComment

	if err := c.ShouldBindUri(&uriParam); err != nil {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusBadRequest,
		})
		return
	}

	isValidatePass, _, validatorError := PTGUvalidator.Validate(uriParam)
	if validatorError != nil {
		logger.Error("[Handler][DeleteMovieCommentHandler()]->Error Validate Data", logger.Field("error", validatorError.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}
	if !isValidatePass {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusBadRequest,
			Error:        "Invalid Parameter",
		})
		return
	}

	databaseTx := database.DB.Begin()
	controllerInstance := controller.NewController(databaseTx)
	defer databaseTx.Rollback()

	deleteComment, err := controllerInstance.DeleteMovieComment(modelController.ParamDeleteMovieComment{
		AccountUUID: c.GetString("ACCOUNT_UUID"),
		CommentUUID: uriParam.CommentUUID,
	})
	if err != nil {
		logger.Error("[Handler][DeleteMovieCommentHandler()]->Error DeleteMovieComment()", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}
	if deleteComment.IsNotFound {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusNotFound,
			Error:        "Comment Not Found",
		})
		return
	}

	databaseTx.Commit()

	utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
		ResponseCode: http.StatusOK,
		Data:         "Delete Comment Success",
	})
}
