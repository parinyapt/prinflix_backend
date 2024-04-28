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
	utilsResponse "github.com/parinyapt/prinflix_backend/utils/response"
)

func AddFavoriteMovieHandler(c *gin.Context) {
	var uriParam modelHandler.UriParamMovieUUIDonly

	if err := c.ShouldBindUri(&uriParam); err != nil {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusBadRequest,
		})
		return
	}

	isValidatePass, _, validatorError := PTGUvalidator.Validate(uriParam)
	if validatorError != nil {
		logger.Error("[Handler][AddFavoriteMovieHandler()]->Error Validate Data", logger.Field("error", validatorError.Error()))
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

	checkFavoriteMovie, err := controllerInstance.CheckFavoriteMovie(modelController.ParamAccountUUIDandMovieUUID{
		AccountUUID: c.GetString("ACCOUNT_UUID"),
		MovieUUID:   uriParam.MovieUUID,
	})
	if err != nil {
		logger.Error("[Handler][AddFavoriteMovieHandler()]->Error CheckFavoriteMovie()", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}
	if !checkFavoriteMovie.IsNotFound {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusBadRequest,
			Error:        "Movie Already Favorited",
		})
		return
	}

	createFavoriteMovieErr := controllerInstance.CreateFavoriteMovie(modelController.ParamAccountUUIDandMovieUUID{
		AccountUUID: c.GetString("ACCOUNT_UUID"),
		MovieUUID:   uriParam.MovieUUID,
	})
	if createFavoriteMovieErr != nil {
		logger.Error("[Handler][AddFavoriteMovieHandler()]->Error CreateFavoriteMovie()", logger.Field("error", createFavoriteMovieErr.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}

	databaseTx.Commit()

	utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
		ResponseCode: http.StatusOK,
		Data:         "Add Favorite Movie Success",
	})
}

func RemoveFavoriteMovieHandler(c *gin.Context) {
	var uriParam modelHandler.UriParamMovieUUIDonly

	if err := c.ShouldBindUri(&uriParam); err != nil {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusBadRequest,
		})
		return
	}

	isValidatePass, _, validatorError := PTGUvalidator.Validate(uriParam)
	if validatorError != nil {
		logger.Error("[Handler][AddFavoriteMovieHandler()]->Error Validate Data", logger.Field("error", validatorError.Error()))
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

	checkFavoriteMovie, err := controllerInstance.CheckFavoriteMovie(modelController.ParamAccountUUIDandMovieUUID{
		AccountUUID: c.GetString("ACCOUNT_UUID"),
		MovieUUID:   uriParam.MovieUUID,
	})
	if err != nil {
		logger.Error("[Handler][AddFavoriteMovieHandler()]->Error CheckFavoriteMovie()", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}
	if checkFavoriteMovie.IsNotFound {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusBadRequest,
			Error:        "Movie Not Favorited",
		})
		return
	}

	deleteFavoriteMovie, err := controllerInstance.DeleteFavoriteMovie(modelController.ParamAccountUUIDandMovieUUID{
		AccountUUID: c.GetString("ACCOUNT_UUID"),
		MovieUUID:   uriParam.MovieUUID,
	})
	if err != nil {
		logger.Error("[Handler][AddFavoriteMovieHandler()]->Error DeleteFavoriteMovie()", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}
	if deleteFavoriteMovie.IsNotFound {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusBadRequest,
			Error:        "Movie Not Favorited",
		})
		return
	}

	databaseTx.Commit()

	utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
		ResponseCode: http.StatusOK,
		Data:         "Remove Favorite Movie Success",
	})
}

func GetFavoriteMovieListHandler(c *gin.Context) {
	controllerInstance := controller.NewController(database.DB)

	var response modelHandler.ResponseGetFavoriteMovieList
	response.ResultData = []modelHandler.ResponseMovieData{}

	getAllFavoriteMovie, err := controllerInstance.GetAllFavoriteMovie(c.GetString("ACCOUNT_UUID"))
	if err != nil {
		logger.Error("[Handler][GetFavoriteMovieListHandler()]->Error GetAllFavoriteMovie()", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}
	if getAllFavoriteMovie.IsNotFound {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusOK,
			Data:         response,
		})
		return
	}

	for _, movie := range getAllFavoriteMovie.Data {
		response.ResultData = append(response.ResultData, modelHandler.ResponseMovieData{
			MovieUUID:            movie.MovieUUID,
			MovieThumbnail:       movie.MovieThumbnail,
			MovieTitle:           movie.MovieTitle,
			MovieDescription:     movie.MovieDescription,
			MovieCategoryID:      movie.MovieCategoryID,
			MovieCategoryName:    movie.MovieCategoryName,
			IsFavorite:           movie.IsFavorite,
			ReviewTotalCount:     movie.ReviewTotalCount,
			ReviewGoodCount:      movie.ReviewGoodCount,
			ReviewFairCount:      movie.ReviewFairCount,
			ReviewBadCount:       movie.ReviewBadCount,
			ReviewGoodPercentage: movie.ReviewGoodPercentage,
			ReviewFairPercentage: movie.ReviewFairPercentage,
			ReviewBadPercentage:  movie.ReviewBadPercentage,
		})
	}

	utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
		ResponseCode: http.StatusOK,
		Data:         response,
	})
}
