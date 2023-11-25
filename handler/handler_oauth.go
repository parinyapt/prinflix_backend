package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/parinyapt/prinflix_backend/controller"
	"github.com/parinyapt/prinflix_backend/database"
	"github.com/parinyapt/prinflix_backend/logger"
	modelController "github.com/parinyapt/prinflix_backend/model/controller"
	modelDatabase "github.com/parinyapt/prinflix_backend/model/database"
	modelHandler "github.com/parinyapt/prinflix_backend/model/handler"
	modelUtils "github.com/parinyapt/prinflix_backend/model/utils"
	utilsResponse "github.com/parinyapt/prinflix_backend/utils/response"
)

func RequestConnectGoogleOAuthHandler(c *gin.Context) {
	var response modelHandler.ResponseRequestConnectOAuth

	databaseTx := database.DB.Begin()
	controllerInstance := controller.NewController(databaseTx)
	defer databaseTx.Rollback()

	checkAccountOAuth, err := controllerInstance.CheckAccountOAuth(modelDatabase.AccountOAuthProviderGoogle ,modelController.ParamCheckAccountOAuth{
		AccountUUID: c.GetString("ACCOUNT_UUID"),
	})
	if err != nil {
		logger.Error("[Handler][RequestConnectGoogleOAuthHandler()]->Error CheckAccountOAuth()", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}
	if !checkAccountOAuth.IsNotFound {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusBadRequest,
			Error: 			"Google OAuth Already Connected",
		})
		return
	}

	err = controllerInstance.DeleteTemporaryCode(modelController.ParamTemporaryCode{
		AccountUUID: c.GetString("ACCOUNT_UUID"),
		Type:        modelDatabase.TemporaryCodeTypeOAuthStateGoogle,
	})
	if err != nil {
		logger.Error("[Handler][RequestConnectGoogleOAuthHandler()]->Error DeleteTemporaryCode()", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}

	createTemporaryCode, err := controllerInstance.CreateTemporaryCode(modelController.ParamTemporaryCode{
		AccountUUID: c.GetString("ACCOUNT_UUID"),
		Type:        modelDatabase.TemporaryCodeTypeOAuthStateGoogle,
	})
	if err != nil {
		logger.Error("[Handler][RequestConnectGoogleOAuthHandler()]->Error CreateTemporaryCode()", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}

	databaseTx.Commit()

	response.AuthURL = controller.GenerateGoogleOAuthURL(createTemporaryCode.CodeUUID.String())

	utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
		ResponseCode: http.StatusOK,
		Data:         response,
	})
}

func RequestDisconnectGoogleOAuthHandler(c *gin.Context) {
	databaseTx := database.DB.Begin()
	controllerInstance := controller.NewController(databaseTx)
	defer databaseTx.Rollback()

	deleteOAuthAccount, err := controllerInstance.DeleteAccountOAuth(modelController.ParamDeleteAccountOAuth{
		AccountUUID: c.GetString("ACCOUNT_UUID"),
		Provider:        modelDatabase.AccountOAuthProviderGoogle,
	})
	if err != nil {
		logger.Error("[Handler][RequestDisconnectGoogleOAuthHandler()]->Error DeleteAccountOAuth()", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}

	if deleteOAuthAccount.IsNotFound {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusNotFound,
			Error: 			"Google OAuth Not Connected",
		})
		return
	}

	databaseTx.Commit()

	utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
		ResponseCode: http.StatusOK,
		Data: 			 "Google OAuth Disconnected",
	})
}

func RequestConnectLineOAuthHandler(c *gin.Context) {
	var response modelHandler.ResponseRequestConnectOAuth

	databaseTx := database.DB.Begin()
	controllerInstance := controller.NewController(databaseTx)
	defer databaseTx.Rollback()

	checkAccountOAuth, err := controllerInstance.CheckAccountOAuth(modelDatabase.AccountOAuthProviderLine ,modelController.ParamCheckAccountOAuth{
		AccountUUID: c.GetString("ACCOUNT_UUID"),
	})
	if err != nil {
		logger.Error("[Handler][RequestConnectLineOAuthHandler()]->Error CheckAccountOAuth()", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}
	if !checkAccountOAuth.IsNotFound {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusBadRequest,
			Error: 			"Line OAuth Already Connected",
		})
		return
	}

	err = controllerInstance.DeleteTemporaryCode(modelController.ParamTemporaryCode{
		AccountUUID: c.GetString("ACCOUNT_UUID"),
		Type:        modelDatabase.TemporaryCodeTypeOAuthStateLine,
	})
	if err != nil {
		logger.Error("[Handler][RequestConnectLineOAuthHandler()]->Error DeleteTemporaryCode()", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}

	createTemporaryCode, err := controllerInstance.CreateTemporaryCode(modelController.ParamTemporaryCode{
		AccountUUID: c.GetString("ACCOUNT_UUID"),
		Type:        modelDatabase.TemporaryCodeTypeOAuthStateLine,
	})
	if err != nil {
		logger.Error("[Handler][RequestConnectLineOAuthHandler()]->Error CreateTemporaryCode()", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}

	databaseTx.Commit()

	response.AuthURL = controller.GenerateLineOAuthURL(createTemporaryCode.CodeUUID.String())

	utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
		ResponseCode: http.StatusOK,
		Data:         response,
	})
}

func RequestDisconnectLineOAuthHandler(c *gin.Context) {
	databaseTx := database.DB.Begin()
	controllerInstance := controller.NewController(databaseTx)
	defer databaseTx.Rollback()

	deleteOAuthAccount, err := controllerInstance.DeleteAccountOAuth(modelController.ParamDeleteAccountOAuth{
		AccountUUID: c.GetString("ACCOUNT_UUID"),
		Provider:        modelDatabase.AccountOAuthProviderLine,
	})
	if err != nil {
		logger.Error("[Handler][RequestDisconnectLineOAuthHandler()]->Error DeleteAccountOAuth()", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}

	if deleteOAuthAccount.IsNotFound {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusNotFound,
			Error: 			"Line OAuth Not Connected",
		})
		return
	}

	databaseTx.Commit()

	utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
		ResponseCode: http.StatusOK,
		Data: 			 "Line OAuth Disconnected",
	})
}