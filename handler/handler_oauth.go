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
	utilsRedirect "github.com/parinyapt/prinflix_backend/utils/redirect"
	utilsResponse "github.com/parinyapt/prinflix_backend/utils/response"
)

func RequestConnectGoogleOAuthHandler(c *gin.Context) {
	var response modelHandler.ResponseRequestConnectOAuth

	databaseTx := database.DB.Begin()
	controllerInstance := controller.NewController(databaseTx)
	defer databaseTx.Rollback()

	checkAccountOAuth, err := controllerInstance.CheckAccountOAuth(modelDatabase.AccountOAuthProviderGoogle, modelController.ParamCheckAccountOAuth{
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
			Error:        "Google OAuth Already Connected",
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
		Provider:    modelDatabase.AccountOAuthProviderGoogle,
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
			Error:        "Google OAuth Not Connected",
		})
		return
	}

	databaseTx.Commit()

	utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
		ResponseCode: http.StatusOK,
		Data:         "Google OAuth Disconnected",
	})
}

func RequestConnectLineOAuthHandler(c *gin.Context) {
	var response modelHandler.ResponseRequestConnectOAuth

	databaseTx := database.DB.Begin()
	controllerInstance := controller.NewController(databaseTx)
	defer databaseTx.Rollback()

	checkAccountOAuth, err := controllerInstance.CheckAccountOAuth(modelDatabase.AccountOAuthProviderLine, modelController.ParamCheckAccountOAuth{
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
			Error:        "Line OAuth Already Connected",
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
		Provider:    modelDatabase.AccountOAuthProviderLine,
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
			Error:        "Line OAuth Not Connected",
		})
		return
	}

	databaseTx.Commit()

	utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
		ResponseCode: http.StatusOK,
		Data:         "Line OAuth Disconnected",
	})
}

func GoogleCallbackHandler(c *gin.Context) {
	var queryParam modelHandler.QueryParamOAuthCallback

	if err := c.ShouldBind(&queryParam); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	isValidatePass, _, validatorError := PTGUvalidator.Validate(queryParam)
	if validatorError != nil {
		logger.Error("[Handler][GoogleCallbackHandler()]->Error Validate()", logger.Field("error", validatorError.Error()))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if !isValidatePass {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	databaseTx := database.DB.Begin()
	controllerInstance := controller.NewController(databaseTx)
	defer databaseTx.Rollback()

	checkTemporaryCode, err := controllerInstance.CheckTemporaryCode(modelController.ParamCheckTemporaryCode{
		CodeUUID: queryParam.State,
		Type:     modelDatabase.TemporaryCodeTypeOAuthStateGoogle,
	})
	if err != nil {
		logger.Error("[Handler][GoogleCallbackHandler()]->Error CheckTemporaryCode()", logger.Field("error", err.Error()))
		c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthConnectRedirectUrl(utilsRedirect.ProviderGoogle, false))
		return
	}
	if checkTemporaryCode.IsNotFound || checkTemporaryCode.IsExpired {
		c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthConnectRedirectUrl(utilsRedirect.ProviderGoogle, false))
		return
	}

	err = controllerInstance.DeleteTemporaryCode(modelController.ParamTemporaryCode{
		AccountUUID: checkTemporaryCode.AccountUUID.String(),
		Type:        modelDatabase.TemporaryCodeTypeOAuthStateGoogle,
	})
	if err != nil {
		logger.Error("[Handler][GoogleCallbackHandler()]->Error DeleteTemporaryCode()", logger.Field("error", err.Error()))
		c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthConnectRedirectUrl(utilsRedirect.ProviderGoogle, false))
		return
	}

	getGoogleOAuthUserInfo, err := controller.GetGoogleOAuthUserInfo(queryParam.Code)
	if err != nil {
		logger.Error("[Handler][GoogleCallbackHandler()]->Error GetGoogleOAuthUserInfo()", logger.Field("error", err.Error()))
		c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthConnectRedirectUrl(utilsRedirect.ProviderGoogle, false))
		return
	}

	err = controllerInstance.CreateAccountOAuth(modelController.ParamCreateAccountOAuth{
		AccountUUID: checkTemporaryCode.AccountUUID.String(),
		Provider:    modelDatabase.AccountOAuthProviderGoogle,
		UserID:      getGoogleOAuthUserInfo.UserID,
		UserName:    getGoogleOAuthUserInfo.Name,
		UserEmail:   getGoogleOAuthUserInfo.Email,
		UserPicture: getGoogleOAuthUserInfo.Picture,
	})
	if err != nil {
		logger.Error("[Handler][GoogleCallbackHandler()]->Error CreateAccountOAuth()", logger.Field("error", err.Error()))
		c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthConnectRedirectUrl(utilsRedirect.ProviderGoogle, false))
		return
	}

	databaseTx.Commit()

	c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthConnectRedirectUrl(utilsRedirect.ProviderGoogle, true))
}

func LineCallbackHandler(c *gin.Context) {
	var queryParam modelHandler.QueryParamOAuthCallback

	if err := c.ShouldBind(&queryParam); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	isValidatePass, _, validatorError := PTGUvalidator.Validate(queryParam)
	if validatorError != nil {
		logger.Error("[Handler][LineCallbackHandler()]->Error Validate()", logger.Field("error", validatorError.Error()))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if !isValidatePass {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	databaseTx := database.DB.Begin()
	controllerInstance := controller.NewController(databaseTx)
	defer databaseTx.Rollback()

	checkTemporaryCode, err := controllerInstance.CheckTemporaryCode(modelController.ParamCheckTemporaryCode{
		CodeUUID: queryParam.State,
		Type:     modelDatabase.TemporaryCodeTypeOAuthStateLine,
	})
	if err != nil {
		logger.Error("[Handler][LineCallbackHandler()]->Error CheckTemporaryCode()", logger.Field("error", err.Error()))
		c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthConnectRedirectUrl(utilsRedirect.ProviderLine, false))
		return
	}
	if checkTemporaryCode.IsNotFound || checkTemporaryCode.IsExpired {
		c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthConnectRedirectUrl(utilsRedirect.ProviderLine, false))
		return
	}

	err = controllerInstance.DeleteTemporaryCode(modelController.ParamTemporaryCode{
		AccountUUID: checkTemporaryCode.AccountUUID.String(),
		Type:        modelDatabase.TemporaryCodeTypeOAuthStateLine,
	})
	if err != nil {
		logger.Error("[Handler][LineCallbackHandler()]->Error DeleteTemporaryCode()", logger.Field("error", err.Error()))
		c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthConnectRedirectUrl(utilsRedirect.ProviderLine, false))
		return
	}

	getLineOAuthUserInfo, err := controller.GetLineOAuthUserInfo(queryParam.Code)
	if err != nil {
		logger.Error("[Handler][LineCallbackHandler()]->Error GetLineOAuthUserInfo()", logger.Field("error", err.Error()))
		c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthConnectRedirectUrl(utilsRedirect.ProviderLine, false))
		return
	}

	err = controllerInstance.CreateAccountOAuth(modelController.ParamCreateAccountOAuth{
		AccountUUID: checkTemporaryCode.AccountUUID.String(),
		Provider:    modelDatabase.AccountOAuthProviderLine,
		UserID:      getLineOAuthUserInfo.UserID,
		UserName:    getLineOAuthUserInfo.Name,
		UserEmail:   getLineOAuthUserInfo.Email,
		UserPicture: getLineOAuthUserInfo.Picture,
	})
	if err != nil {
		logger.Error("[Handler][LineCallbackHandler()]->Error CreateAccountOAuth()", logger.Field("error", err.Error()))
		c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthConnectRedirectUrl(utilsRedirect.ProviderLine, false))
		return
	}

	databaseTx.Commit()

	c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthConnectRedirectUrl(utilsRedirect.ProviderLine, true))
}
