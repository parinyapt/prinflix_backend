package handler

import (
	"encoding/base64"
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

func RequestWatchMovieHandler(c *gin.Context) {
	var uriParam modelHandler.UriParamMovieUUIDonly
	var response modelHandler.ResponseRequestWatchMovie

	if err := c.ShouldBindUri(&uriParam); err != nil {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusBadRequest,
		})
		return
	}

	isValidatePass, _, validatorError := PTGUvalidator.Validate(uriParam)
	if validatorError != nil {
		logger.Error("[Handler][RequestWatchMovieHandler()]->Error Validate Data", logger.Field("error", validatorError.Error()))
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

	movieDetail, err := controllerInstance.GetMovieDetail(modelController.ParamAccountUUIDandMovieUUID{
		AccountUUID: c.GetString("ACCOUNT_UUID"),
		MovieUUID:   uriParam.MovieUUID,
	})
	if err != nil {
		logger.Error("[Handler][RequestWatchMovieHandler()]->Error GetMovieDetail()", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}
	if movieDetail.IsNotFound {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusNotFound,
			Error:        "Movie Not Found",
		})
		return
	}

	err = controllerInstance.DeleteAllWatchSessionByAccountUUID(c.GetString("ACCOUNT_UUID"))
	if err != nil {
		logger.Error("[Handler][RequestWatchMovieHandler()]->Error DeleteAllWatchSessionByAccountUUID()", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}

	createWatchSession, err := controllerInstance.CreateWatchSession(modelController.ParamAccountUUIDandMovieUUID{
		AccountUUID: c.GetString("ACCOUNT_UUID"),
		MovieUUID:   uriParam.MovieUUID,
	})
	if err != nil {
		logger.Error("[Handler][RequestWatchMovieHandler()]->Error CreateWatchSession()", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}

	updateWatchHistory, err := controllerInstance.ClearWatchHistoryStart(modelController.ParamUpdateWatchHistory{
		AccountUUID: c.GetString("ACCOUNT_UUID"),
		MovieUUID:   uriParam.MovieUUID,
	})
	if err != nil {
		logger.Error("[Handler][RequestWatchMovieHandler()]->Error ClearWatchHistoryStart()", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}
	if updateWatchHistory.IsNotFound {
		err := controllerInstance.CreateWatchHistory(modelController.ParamAccountUUIDandMovieUUID{
			AccountUUID: c.GetString("ACCOUNT_UUID"),
			MovieUUID:   uriParam.MovieUUID,
		})
		if err != nil {
			logger.Error("[Handler][RequestWatchMovieHandler()]->Error CreateWatchHistory()", logger.Field("error", err.Error()))
			utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
				ResponseCode: http.StatusInternalServerError,
			})
			return
		}
	}

	generateWatchSessionToken, err := controller.GenerateWatchSessionToken(modelController.ParamGenerateWatchSessionToken{
		SessionUUID: createWatchSession.SessionUUID.String(),
		ExpiredAt:   createWatchSession.ExpiredAt,
	})
	if err != nil {
		logger.Error("[Handler][LoginHandler()]->Error Generate Access Token", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}

	databaseTx.Commit()

	cookieMaxAge := int(controller.WatchSessionExpiredIn.Seconds())
	c.SetCookie("prinflix_session_token", generateWatchSessionToken.WatchSessionToken, cookieMaxAge, "/", "prinpt.com", true, true)
	response.WatchSessionToken = base64.URLEncoding.EncodeToString([]byte(generateWatchSessionToken.WatchSessionToken))
	utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
		ResponseCode: http.StatusOK,
		Data:         response,
	})
}

func RequestPauseMovieHandler(c *gin.Context) {
	var request modelHandler.RequestPauseMovie

	if err := c.ShouldBindJSON(&request); err != nil {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusBadRequest,
		})
		return
	}

	isValidatePass, errorFieldList, validatorError := PTGUvalidator.Validate(request)
	if validatorError != nil {
		logger.Error("[Handler][RequestPauseMovieHandler()]->Error Validate Data", logger.Field("error", validatorError.Error()))
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

	updateWatchHistory, err := controllerInstance.ClearWatchHistoryPause(modelController.ParamUpdateWatchHistory{
		AccountUUID: c.GetString("ACCOUNT_UUID"),
		MovieUUID:   c.GetString("WATCHSESSION_MOVIE_UUID"),
		TimeStamp:   request.TimeStamp,
	})
	if err != nil {
		logger.Error("[Handler][RequestPauseMovieHandler()]->Error ClearWatchHistoryPause()", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}
	if updateWatchHistory.IsNotFound {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusBadRequest,
			Error:        "Invalid Session",
		})
		return
	}

	databaseTx.Commit()

	utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
		ResponseCode: http.StatusOK,
		Data:         "Update Watch History Success",
	})
}

func RequestEndMovieHandler(c *gin.Context) {
	databaseTx := database.DB.Begin()
	controllerInstance := controller.NewController(databaseTx)
	defer databaseTx.Rollback()

	err := controllerInstance.DeleteAllWatchSessionByAccountUUID(c.GetString("ACCOUNT_UUID"))
	if err != nil {
		logger.Error("[Handler][RequestEndMovieHandler()]->Error DeleteAllWatchSessionByAccountUUID()", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}

	updateWatchHistory, err := controllerInstance.ClearWatchHistoryEnd(modelController.ParamUpdateWatchHistory{
		AccountUUID: c.GetString("ACCOUNT_UUID"),
		MovieUUID:   c.GetString("WATCHSESSION_MOVIE_UUID"),
	})
	if err != nil {
		logger.Error("[Handler][RequestEndMovieHandler()]->Error ClearWatchHistoryEnd()", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}
	if updateWatchHistory.IsNotFound {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusBadRequest,
			Error:        "Invalid Session",
		})
		return
	}

	databaseTx.Commit()

	c.SetCookie("prinflix_session_token", "", -1, "/", "prinpt.com", true, true)
	utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
		ResponseCode: http.StatusOK,
		Data:         "Update Watch History Success",
	})
}
