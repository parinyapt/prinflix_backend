package handler

import (
	"net/http"

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

func GetMovieListHandler(c *gin.Context) {
	var queryParam modelHandler.QueryParamGetMovieList

	if err := c.ShouldBind(&queryParam); err != nil {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusBadRequest,
		})
		return
	}

	isValidatePass, _, validatorError := PTGUvalidator.Validate(queryParam)
	if validatorError != nil {
		logger.Error("[Handler][GetTipsByMatchUUID()]->Error Validate Data", logger.Field("error", validatorError.Error()))
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

	if queryParam.Pagination.Page == 0 {
		queryParam.Pagination.Page = 1
	}

	if queryParam.Pagination.Limit == 0 {
		queryParam.Pagination.Limit = 10
	}

	queryParam.Pagination.SortField = repository.FetchManyMovieSortFieldMovieTitle

	if len(queryParam.Pagination.SortOrderBy) == 0 {
		queryParam.Pagination.SortOrderBy = repository.SortOrderByAsc
	}

	controllerInstance := controller.NewController(database.DB)

	getManyMovie, err := controllerInstance.GetAllMovie(modelController.ParamGetAllMovie{
		AccountUUID: c.GetString("ACCOUNT_UUID"),
		CategoryID:  queryParam.CategoryID,
		SearchQuery: queryParam.Keyword,
		Pagination: modelController.ParamPagination{
			Page:        queryParam.Pagination.Page,
			Limit:       queryParam.Pagination.Limit,
			SortField:   queryParam.Pagination.SortField,
			SortOrderBy: queryParam.Pagination.SortOrderBy,
		},
	})
	if err != nil {
		logger.Error("[Handler][GetMovieListHandler()]->Error GetAllMovie()", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}
	if getManyMovie.IsNotFound {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusNotFound,
			Error:        "Movie Not Found",
		})
		return
	}

	var response modelHandler.ResponseGetMovieList
	response.ResultPagination.TotalData = getManyMovie.Pagination.TotalData
	response.ResultPagination.TotalPage = getManyMovie.Pagination.TotalPage
	response.ResultPagination.CurrentPage = getManyMovie.Pagination.Page
	response.ResultPagination.CurrentLimit = getManyMovie.Pagination.Limit

	for _, movie := range getManyMovie.Data {
		response.ResultData = append(response.ResultData, modelHandler.ResponseGetMovieListData{
			MovieUUID:         movie.MovieUUID,
			MovieThumbnail:    movie.MovieThumbnail,
			MovieTitle:        movie.MovieTitle,
			MovieDescription:  movie.MovieDescription,
			MovieCategoryID:   movie.MovieCategoryID,
			MovieCategoryName: movie.MovieCategoryName,
			IsFavorite:        movie.IsFavorite,
		})
	}

	utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
		ResponseCode: http.StatusOK,
		Error:        response,
	})
}

func GetMovieCategoryListHandler(c *gin.Context) {
	var response []modelHandler.ResponseGetMovieCategoryList

	controllerInstance := controller.NewController(database.DB)

	getAllMovieCategory, err := controllerInstance.GetAllMovieCategory()
	if err != nil {
		logger.Error("[Handler][GetMovieCategoryListHandler()]->Error GetAllMovieCategory()", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}
	if getAllMovieCategory.IsNotFound {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusNotFound,
			Error:        "Movie Category Not Found",
		})
		return
	}

	for _, movieCategory := range getAllMovieCategory.Data {
		response = append(response, modelHandler.ResponseGetMovieCategoryList{
			CategoryID:   movieCategory.CategoryID,
			CategoryName: movieCategory.CategoryName,
		})
	}

	utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
		ResponseCode: http.StatusOK,
		Data:         response,
	})
}
