package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	PTGUvalidator "github.com/parinyapt/golang_utils/validator/v1"
	"github.com/parinyapt/prinflix_backend/controller"
	"github.com/parinyapt/prinflix_backend/database"
	"github.com/parinyapt/prinflix_backend/logger"
	modelController "github.com/parinyapt/prinflix_backend/model/controller"
	modelDatabase "github.com/parinyapt/prinflix_backend/model/database"
	modelHandler "github.com/parinyapt/prinflix_backend/model/handler"
	modelUtils "github.com/parinyapt/prinflix_backend/model/utils"
	utilsResponse "github.com/parinyapt/prinflix_backend/utils/response"
)

func ReviewMovieHandler(c *gin.Context) {
	var uriParam modelHandler.UriParamMovieUUIDonly
	if err := c.ShouldBindUri(&uriParam); err != nil {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusBadRequest,
		})
		return
	}

	isValidatePass, _, validatorError := PTGUvalidator.Validate(uriParam)
	if validatorError != nil {
		logger.Error("[Handler][ReviewMovieHandler()]->Error Validate Data", logger.Field("error", validatorError.Error()))
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

	var request modelHandler.RequestReviewMovie
	if err := c.ShouldBindJSON(&request); err != nil {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusBadRequest,
		})
		return
	}

	isValidatePass, errorFieldList, validatorError := PTGUvalidator.Validate(request)
	if validatorError != nil {
		logger.Error("[Handler][ReviewMovieHandler()]->Error Validate Data", logger.Field("error", validatorError.Error()))
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

	getMovieDetail, err := controllerInstance.GetMovieDetail(modelController.ParamAccountUUIDandMovieUUID{
		AccountUUID: c.GetString("ACCOUNT_UUID"),
		MovieUUID:   uriParam.MovieUUID,
	})
	if err != nil {
		logger.Error("[Handler][ReviewMovieHandler()]->Error GetMovieDetail()", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}
	if getMovieDetail.IsNotFound {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusNotFound,
			Error:        "Movie Not Found",
		})
		return
	}

	if len(request.Rating) == 0 {
		err := controllerInstance.DeleteReviewMovie(modelController.ParamDeleteReviewMovie{
			AccountUUID: c.GetString("ACCOUNT_UUID"),
			MovieUUID:   getMovieDetail.MovieUUID.String(),
		})
		if err != nil {
			logger.Error("[Handler][ReviewMovieHandler()]->Error DeleteReviewMovie()", logger.Field("error", err.Error()))
			utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
				ResponseCode: http.StatusInternalServerError,
			})
			return
		}
	} else {
		var ratingNumber uint
		switch request.Rating {
		case "good":
			ratingNumber = modelDatabase.ReviewRatingGood
		case "fair":
			ratingNumber = modelDatabase.ReviewRatingFair
		case "bad":
			ratingNumber = modelDatabase.ReviewRatingBad
		default:
			utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
				ResponseCode: http.StatusBadRequest,
			})
			return
		}

		err := controllerInstance.CreateUpdateReviewMovie(modelController.ParamCreateUpdateReviewMovie{
			AccountUUID:  c.GetString("ACCOUNT_UUID"),
			MovieUUID:    getMovieDetail.MovieUUID.String(),
			ReviewRating: ratingNumber,
		})
		if err != nil {
			logger.Error("[Handler][ReviewMovieHandler()]->Error CreateUpdateReviewMovie()", logger.Field("error", err.Error()))
			utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
				ResponseCode: http.StatusInternalServerError,
			})
			return
		}
	}

	databaseTx.Commit()

	utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
		ResponseCode: http.StatusOK,
		Data:         "Review Movie Success",
	})
}
