package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	"github.com/sethvargo/go-password/password"
)

// Connect Google OAuth
func RequestConnectGoogleOAuthHandler(c *gin.Context) {
	var queryParam modelHandler.QueryParamOAuthConnect

	if err := c.ShouldBindQuery(&queryParam); err != nil {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusBadRequest,
		})
		return
	}

	isValidatePass, _, validatorError := PTGUvalidator.Validate(queryParam)
	if validatorError != nil {
		logger.Error("[Handler][RequestConnectGoogleOAuthHandler()]->Error Validate Data", logger.Field("error", validatorError.Error()))
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
	err = controllerInstance.DeleteTemporaryCode(modelController.ParamTemporaryCode{
		AccountUUID: c.GetString("ACCOUNT_UUID"),
		Type:        modelDatabase.TemporaryCodeTypeAppOAuthStateGoogle,
	})
	if err != nil {
		logger.Error("[Handler][RequestConnectGoogleOAuthHandler()]->Error DeleteTemporaryCode()", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}

	var typeConnect string
	if queryParam.IsApplication {
		typeConnect = modelDatabase.TemporaryCodeTypeAppOAuthStateGoogle
	} else {
		typeConnect = modelDatabase.TemporaryCodeTypeOAuthStateGoogle
	}
	createTemporaryCode, err := controllerInstance.CreateTemporaryCode(modelController.ParamTemporaryCode{
		AccountUUID: c.GetString("ACCOUNT_UUID"),
		Type:        typeConnect,
	})
	if err != nil {
		logger.Error("[Handler][RequestConnectGoogleOAuthHandler()]->Error CreateTemporaryCode()", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}

	databaseTx.Commit()

	response.AuthURL = controller.GenerateGoogleOAuthURL(createTemporaryCode.CodeUUID.String(), 1)

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

func GoogleCallbackHandler(c *gin.Context) {
	var queryParam modelHandler.QueryParamOAuthCallback

	if err := c.ShouldBindQuery(&queryParam); err != nil {
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

	var checkType string = modelDatabase.TemporaryCodeTypeOAuthStateGoogle
	var checkTemporaryCode modelController.ReturnCheckTemporaryCode
	checkTemporaryCode, err := controllerInstance.CheckTemporaryCode(modelController.ParamCheckTemporaryCode{
		CodeUUID: queryParam.State,
		Type:     checkType,
	})
	if err != nil {
		logger.Error("[Handler][GoogleCallbackHandler()]->Error CheckTemporaryCode()", logger.Field("error", err.Error()))
		c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthConnectRedirectUrl(utilsRedirect.ProviderGoogle, false))
		return
	}
	if checkTemporaryCode.IsNotFound {
		checkType = modelDatabase.TemporaryCodeTypeAppOAuthStateGoogle
		checkTemporaryCode, err = controllerInstance.CheckTemporaryCode(modelController.ParamCheckTemporaryCode{
			CodeUUID: queryParam.State,
			Type:     checkType,
		})
		if err != nil {
			logger.Error("[Handler][GoogleCallbackHandler()]->Error CheckTemporaryCode()", logger.Field("error", err.Error()))
			c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthConnectRedirectUrl(utilsRedirect.ProviderGoogle, false))
			return
		}
	}
	if checkTemporaryCode.IsExpired {
		c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthConnectRedirectUrl(utilsRedirect.ProviderGoogle, false))
		return
	}

	err = controllerInstance.DeleteTemporaryCode(modelController.ParamTemporaryCode{
		AccountUUID: checkTemporaryCode.AccountUUID.String(),
		Type:        checkType,
	})
	if err != nil {
		logger.Error("[Handler][GoogleCallbackHandler()]->Error DeleteTemporaryCode()", logger.Field("error", err.Error()))
		c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthConnectRedirectUrl(utilsRedirect.ProviderGoogle, false))
		return
	}

	getGoogleOAuthUserInfo, err := controller.GetGoogleOAuthUserInfo(queryParam.Code, 1)
	if err != nil {
		logger.Error("[Handler][GoogleCallbackHandler()]->Error GetGoogleOAuthUserInfo()", logger.Field("error", err.Error()))
		c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthConnectRedirectUrl(utilsRedirect.ProviderGoogle, false))
		return
	}

	checkAccountOAuth, err := controllerInstance.CheckAccountOAuth(modelDatabase.AccountOAuthProviderGoogle, modelController.ParamCheckAccountOAuth{
		UserID: getGoogleOAuthUserInfo.UserID,
	})
	if err != nil {
		logger.Error("[Handler][GoogleCallbackHandler()]->Error CheckAccountOAuth()", logger.Field("error", err.Error()))
		c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthLoginRedirectUrl(utilsRedirect.ProviderGoogle, false, ""))
		return
	}

	if !checkAccountOAuth.IsNotFound {
		if checkType == modelDatabase.TemporaryCodeTypeAppOAuthStateGoogle {
			c.Redirect(http.StatusFound, utilsRedirect.GenerateAppOAuthConnectRedirectUrl(utilsRedirect.ProviderGoogle, false))
		} else {
			c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthConnectRedirectUrl(utilsRedirect.ProviderGoogle, false))
		}
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

	if checkType == modelDatabase.TemporaryCodeTypeAppOAuthStateGoogle {
		c.Redirect(http.StatusFound, utilsRedirect.GenerateAppOAuthConnectRedirectUrl(utilsRedirect.ProviderGoogle, true))
	} else {
		c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthConnectRedirectUrl(utilsRedirect.ProviderGoogle, true))
	}
}

// Connect Line OAuth
func RequestConnectLineOAuthHandler(c *gin.Context) {
	var queryParam modelHandler.QueryParamOAuthConnect

	if err := c.ShouldBindQuery(&queryParam); err != nil {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusBadRequest,
		})
		return
	}

	isValidatePass, _, validatorError := PTGUvalidator.Validate(queryParam)
	if validatorError != nil {
		logger.Error("[Handler][RequestConnectLineOAuthHandler()]->Error Validate Data", logger.Field("error", validatorError.Error()))
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
	err = controllerInstance.DeleteTemporaryCode(modelController.ParamTemporaryCode{
		AccountUUID: c.GetString("ACCOUNT_UUID"),
		Type:        modelDatabase.TemporaryCodeTypeAppOAuthStateLine,
	})
	if err != nil {
		logger.Error("[Handler][RequestConnectLineOAuthHandler()]->Error DeleteTemporaryCode()", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}

	var typeConnect string
	if queryParam.IsApplication {
		typeConnect = modelDatabase.TemporaryCodeTypeAppOAuthStateLine
	} else {
		typeConnect = modelDatabase.TemporaryCodeTypeOAuthStateLine
	}
	createTemporaryCode, err := controllerInstance.CreateTemporaryCode(modelController.ParamTemporaryCode{
		AccountUUID: c.GetString("ACCOUNT_UUID"),
		Type:        typeConnect,
	})
	if err != nil {
		logger.Error("[Handler][RequestConnectLineOAuthHandler()]->Error CreateTemporaryCode()", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}

	databaseTx.Commit()

	response.AuthURL = controller.GenerateLineOAuthURL(createTemporaryCode.CodeUUID.String(), 1)

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

func LineCallbackHandler(c *gin.Context) {
	var queryParam modelHandler.QueryParamOAuthCallback

	if err := c.ShouldBindQuery(&queryParam); err != nil {
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

	var checkType string = modelDatabase.TemporaryCodeTypeOAuthStateLine
	var checkTemporaryCode modelController.ReturnCheckTemporaryCode
	checkTemporaryCode, err := controllerInstance.CheckTemporaryCode(modelController.ParamCheckTemporaryCode{
		CodeUUID: queryParam.State,
		Type:     checkType,
	})
	if err != nil {
		logger.Error("[Handler][LineCallbackHandler()]->Error CheckTemporaryCode()", logger.Field("error", err.Error()))
		c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthConnectRedirectUrl(utilsRedirect.ProviderLine, false))
		return
	}
	if checkTemporaryCode.IsNotFound {
		checkType = modelDatabase.TemporaryCodeTypeAppOAuthStateLine
		checkTemporaryCode, err = controllerInstance.CheckTemporaryCode(modelController.ParamCheckTemporaryCode{
			CodeUUID: queryParam.State,
			Type:     checkType,
		})
		if err != nil {
			logger.Error("[Handler][LineCallbackHandler()]->Error CheckTemporaryCode()", logger.Field("error", err.Error()))
			c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthConnectRedirectUrl(utilsRedirect.ProviderLine, false))
			return
		}
	}
	if checkTemporaryCode.IsExpired {
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

	getLineOAuthUserInfo, err := controller.GetLineOAuthUserInfo(queryParam.Code, 1)
	if err != nil {
		logger.Error("[Handler][LineCallbackHandler()]->Error GetLineOAuthUserInfo()", logger.Field("error", err.Error()))
		c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthConnectRedirectUrl(utilsRedirect.ProviderLine, false))
		return
	}

	checkAccountOAuth, err := controllerInstance.CheckAccountOAuth(modelDatabase.AccountOAuthProviderLine, modelController.ParamCheckAccountOAuth{
		UserID: getLineOAuthUserInfo.UserID,
	})
	if err != nil {
		logger.Error("[Handler][LineCallbackHandler()]->Error CheckAccountOAuth()", logger.Field("error", err.Error()))
		c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthLoginRedirectUrl(utilsRedirect.ProviderLine, false, ""))
		return
	}

	if !checkAccountOAuth.IsNotFound {
		if checkType == modelDatabase.TemporaryCodeTypeAppOAuthStateLine {
			c.Redirect(http.StatusFound, utilsRedirect.GenerateAppOAuthConnectRedirectUrl(utilsRedirect.ProviderLine, false))
		} else {
			c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthConnectRedirectUrl(utilsRedirect.ProviderLine, false))
		}
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

	if checkType == modelDatabase.TemporaryCodeTypeAppOAuthStateLine {
		c.Redirect(http.StatusFound, utilsRedirect.GenerateAppOAuthConnectRedirectUrl(utilsRedirect.ProviderLine, true))
	} else {
		c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthConnectRedirectUrl(utilsRedirect.ProviderLine, true))
	}
}

// Google OAuth Login
func GoogleLoginV2Handler(c *gin.Context) {
	databaseTx := database.DB.Begin()
	controllerInstance := controller.NewController(databaseTx)
	defer databaseTx.Rollback()

	createOauthState, err := controllerInstance.CreateOauthState(modelDatabase.OauthStateProviderGoogle)
	if err != nil {
		logger.Error("[Handler][GoogleLoginV2Handler()]->Error CreateOauthState()", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}

	databaseTx.Commit()

	authURL := controller.GenerateGoogleOAuthURL(createOauthState.StateUUID.String(), 2)

	c.Redirect(http.StatusFound, authURL)
}

func GoogleLoginCallbackV2Handler(c *gin.Context) {
	var queryParam modelHandler.QueryParamOAuthCallback

	if err := c.ShouldBindQuery(&queryParam); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	isValidatePass, _, validatorError := PTGUvalidator.Validate(queryParam)
	if validatorError != nil {
		logger.Error("[Handler][GoogleCallbackV2Handler()]->Error Validate()", logger.Field("error", validatorError.Error()))
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

	checkOauthState, err := controllerInstance.CheckOauthState(modelController.ParamOauthState{
		StateUUID: queryParam.State,
		Provider:  modelDatabase.OauthStateProviderGoogle,
	})
	if err != nil {
		logger.Error("[Handler][GoogleCallbackV2Handler()]->Error CheckOauthState()", logger.Field("error", err.Error()))
		c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthLoginRedirectUrl(utilsRedirect.ProviderGoogle, false, ""))
		return
	}
	if checkOauthState.IsNotFound || checkOauthState.IsExpired {
		c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthLoginRedirectUrl(utilsRedirect.ProviderGoogle, false, ""))
		return
	}

	err = controllerInstance.DeleteOauthState(modelController.ParamOauthState{
		StateUUID: queryParam.State,
		Provider:  modelDatabase.OauthStateProviderGoogle,
	})
	if err != nil {
		logger.Error("[Handler][GoogleCallbackV2Handler()]->Error DeleteOauthState()", logger.Field("error", err.Error()))
		c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthLoginRedirectUrl(utilsRedirect.ProviderGoogle, false, ""))
		return
	}

	googleOAuthUserInfo, err := controller.GetGoogleOAuthUserInfo(queryParam.Code, 2)
	if err != nil {
		logger.Error("[Handler][GoogleCallbackV2Handler()]->Error GetGoogleOAuthUserInfo()", logger.Field("error", err.Error()))
		c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthLoginRedirectUrl(utilsRedirect.ProviderGoogle, false, ""))
		return
	}

	checkAccountOAuth, err := controllerInstance.CheckAccountOAuth(modelDatabase.AccountOAuthProviderGoogle, modelController.ParamCheckAccountOAuth{
		UserID: googleOAuthUserInfo.UserID,
	})
	if err != nil {
		logger.Error("[Handler][GoogleCallbackV2Handler()]->Error CheckAccountOAuth()", logger.Field("error", err.Error()))
		c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthLoginRedirectUrl(utilsRedirect.ProviderGoogle, false, ""))
		return
	}

	var accountUUID uuid.UUID = checkAccountOAuth.AccountUUID

	if checkAccountOAuth.IsNotFound {
		password, err := password.Generate(64, 10, 10, false, false)
		if err != nil {
			logger.Error("[Handler][GoogleCallbackV2Handler()]->Error Generate Password", logger.Field("error", err.Error()))
			c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthLoginRedirectUrl(utilsRedirect.ProviderGoogle, false, ""))
			return
		}
		createAccount, err := controllerInstance.CreateAccount(modelController.ParamCreateAccount{
			Name:               googleOAuthUserInfo.Name,
			Email:              googleOAuthUserInfo.Email,
			Password:           password,
			EmailVerifyApprove: true,
		})
		if err != nil {
			logger.Error("[Handler][GoogleCallbackV2Handler()]->Error Create Account", logger.Field("error", err.Error()))
			c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthLoginRedirectUrl(utilsRedirect.ProviderGoogle, false, ""))
			return
		}

		accountUUID = createAccount.UUID

		err = controllerInstance.CreateAccountOAuth(modelController.ParamCreateAccountOAuth{
			AccountUUID: accountUUID.String(),
			Provider:    modelDatabase.AccountOAuthProviderGoogle,
			UserID:      googleOAuthUserInfo.UserID,
			UserName:    googleOAuthUserInfo.Name,
			UserEmail:   googleOAuthUserInfo.Email,
			UserPicture: googleOAuthUserInfo.Picture,
		})
		if err != nil {
			logger.Error("[Handler][GoogleCallbackV2Handler()]->Error CreateAccountOAuth()", logger.Field("error", err.Error()))
			c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthLoginRedirectUrl(utilsRedirect.ProviderGoogle, false, ""))
			return
		}
	}

	createTemporaryCode, err := controllerInstance.CreateTemporaryCode(modelController.ParamTemporaryCode{
		AccountUUID: accountUUID.String(),
		Type:        modelDatabase.TemporaryCodeTypeAuthTokenCode,
	})
	if err != nil {
		logger.Error("[Handler][GoogleCallbackV2Handler()]->Error CreateTemporaryCode()", logger.Field("error", err.Error()))
		c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthLoginRedirectUrl(utilsRedirect.ProviderGoogle, false, ""))
		return
	}

	codeUUIDEncryptBase64, err := controller.EncryptTemporaryCode(createTemporaryCode.CodeUUID.String())
	if err != nil {
		logger.Error("[Handler][GoogleCallbackV2Handler()]->Error EncryptTemporaryCode()", logger.Field("error", err.Error()))
		c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthLoginRedirectUrl(utilsRedirect.ProviderGoogle, false, ""))
		return
	}

	databaseTx.Commit()

	c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthLoginRedirectUrl(utilsRedirect.ProviderGoogle, true, codeUUIDEncryptBase64))
}

// Line OAuth Login
func LineLoginV2Handler(c *gin.Context) {
	databaseTx := database.DB.Begin()
	controllerInstance := controller.NewController(databaseTx)
	defer databaseTx.Rollback()

	createOauthState, err := controllerInstance.CreateOauthState(modelDatabase.OauthStateProviderLine)
	if err != nil {
		logger.Error("[Handler][LineLoginV2Handler()]->Error CreateOauthState()", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}

	databaseTx.Commit()

	authURL := controller.GenerateLineOAuthURL(createOauthState.StateUUID.String(), 2)

	c.Redirect(http.StatusFound, authURL)
}

func LineLoginCallbackV2Handler(c *gin.Context) {
	var queryParam modelHandler.QueryParamOAuthCallback

	if err := c.ShouldBindQuery(&queryParam); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	isValidatePass, _, validatorError := PTGUvalidator.Validate(queryParam)
	if validatorError != nil {
		logger.Error("[Handler][LineCallbackV2Handler()]->Error Validate()", logger.Field("error", validatorError.Error()))
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

	checkOauthState, err := controllerInstance.CheckOauthState(modelController.ParamOauthState{
		StateUUID: queryParam.State,
		Provider:  modelDatabase.OauthStateProviderLine,
	})
	if err != nil {
		logger.Error("[Handler][LineCallbackV2Handler()]->Error CheckOauthState()", logger.Field("error", err.Error()))
		c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthLoginRedirectUrl(utilsRedirect.ProviderLine, false, ""))
		return
	}
	if checkOauthState.IsNotFound || checkOauthState.IsExpired {
		c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthLoginRedirectUrl(utilsRedirect.ProviderLine, false, ""))
		return
	}

	err = controllerInstance.DeleteOauthState(modelController.ParamOauthState{
		StateUUID: queryParam.State,
		Provider:  modelDatabase.OauthStateProviderLine,
	})
	if err != nil {
		logger.Error("[Handler][LineCallbackV2Handler()]->Error DeleteOauthState()", logger.Field("error", err.Error()))
		c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthLoginRedirectUrl(utilsRedirect.ProviderLine, false, ""))
		return
	}

	lineOAuthUserInfo, err := controller.GetLineOAuthUserInfo(queryParam.Code, 2)
	if err != nil {
		logger.Error("[Handler][LineCallbackV2Handler()]->Error GetLineOAuthUserInfo()", logger.Field("error", err.Error()))
		c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthConnectRedirectUrl(utilsRedirect.ProviderLine, false))
		return
	}

	checkAccountOAuth, err := controllerInstance.CheckAccountOAuth(modelDatabase.AccountOAuthProviderLine, modelController.ParamCheckAccountOAuth{
		UserID: lineOAuthUserInfo.UserID,
	})
	if err != nil {
		logger.Error("[Handler][LineCallbackV2Handler()]->Error CheckAccountOAuth()", logger.Field("error", err.Error()))
		c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthLoginRedirectUrl(utilsRedirect.ProviderLine, false, ""))
		return
	}

	var accountUUID uuid.UUID = checkAccountOAuth.AccountUUID

	if checkAccountOAuth.IsNotFound {
		password, err := password.Generate(64, 10, 10, false, false)
		if err != nil {
			logger.Error("[Handler][LineCallbackV2Handler()]->Error Generate Password", logger.Field("error", err.Error()))
			c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthLoginRedirectUrl(utilsRedirect.ProviderLine, false, ""))
			return
		}
		createAccount, err := controllerInstance.CreateAccount(modelController.ParamCreateAccount{
			Name:               lineOAuthUserInfo.Name,
			Email:              lineOAuthUserInfo.Email,
			Password:           password,
			EmailVerifyApprove: true,
		})
		if err != nil {
			logger.Error("[Handler][LineCallbackV2Handler()]->Error Create Account", logger.Field("error", err.Error()))
			c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthLoginRedirectUrl(utilsRedirect.ProviderLine, false, ""))
			return
		}

		accountUUID = createAccount.UUID

		err = controllerInstance.CreateAccountOAuth(modelController.ParamCreateAccountOAuth{
			AccountUUID: accountUUID.String(),
			Provider:    modelDatabase.AccountOAuthProviderLine,
			UserID:      lineOAuthUserInfo.UserID,
			UserName:    lineOAuthUserInfo.Name,
			UserEmail:   lineOAuthUserInfo.Email,
			UserPicture: lineOAuthUserInfo.Picture,
		})
		if err != nil {
			logger.Error("[Handler][LineCallbackV2Handler()]->Error CreateAccountOAuth()", logger.Field("error", err.Error()))
			c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthLoginRedirectUrl(utilsRedirect.ProviderLine, false, ""))
			return
		}
	}

	createTemporaryCode, err := controllerInstance.CreateTemporaryCode(modelController.ParamTemporaryCode{
		AccountUUID: accountUUID.String(),
		Type:        modelDatabase.TemporaryCodeTypeAuthTokenCode,
	})
	if err != nil {
		logger.Error("[Handler][LineCallbackV2Handler()]->Error CreateTemporaryCode()", logger.Field("error", err.Error()))
		c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthLoginRedirectUrl(utilsRedirect.ProviderLine, false, ""))
		return
	}

	codeUUIDEncryptBase64, err := controller.EncryptTemporaryCode(createTemporaryCode.CodeUUID.String())
	if err != nil {
		logger.Error("[Handler][LineCallbackV2Handler()]->Error EncryptTemporaryCode()", logger.Field("error", err.Error()))
		c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthLoginRedirectUrl(utilsRedirect.ProviderLine, false, ""))
		return
	}

	databaseTx.Commit()

	c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthLoginRedirectUrl(utilsRedirect.ProviderLine, true, codeUUIDEncryptBase64))
}

// Apple OAuth Login
func AppleLoginV2Handler(c *gin.Context) {
	databaseTx := database.DB.Begin()
	controllerInstance := controller.NewController(databaseTx)
	defer databaseTx.Rollback()

	createOauthState, err := controllerInstance.CreateOauthState(modelDatabase.OauthStateProviderApple)
	if err != nil {
		logger.Error("[Handler][AppleLoginV2Handler()]->Error CreateOauthState()", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}

	databaseTx.Commit()

	authURL := controller.GenerateAppleOAuthURL(createOauthState.StateUUID.String(), 2)

	c.Redirect(http.StatusFound, authURL)
}

func AppleLoginCallbackV2Handler(c *gin.Context) {
	var request modelHandler.RequestAppleCallback

	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	isValidatePass, _, validatorError := PTGUvalidator.Validate(request)
	if validatorError != nil {
		logger.Error("[Handler][AppleCallbackV2Handler()]->Error Validate Data", logger.Field("error", validatorError.Error()))
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

	checkOauthState, err := controllerInstance.CheckOauthState(modelController.ParamOauthState{
		StateUUID: request.State,
		Provider:  modelDatabase.OauthStateProviderApple,
	})
	if err != nil {
		logger.Error("[Handler][AppleCallbackV2Handler()]->Error CheckOauthState()", logger.Field("error", err.Error()))
		c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthLoginRedirectUrl(utilsRedirect.ProviderApple, false, ""))
		return
	}
	if checkOauthState.IsNotFound || checkOauthState.IsExpired {
		c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthLoginRedirectUrl(utilsRedirect.ProviderApple, false, ""))
		return
	}

	err = controllerInstance.DeleteOauthState(modelController.ParamOauthState{
		StateUUID: request.State,
		Provider:  modelDatabase.OauthStateProviderApple,
	})
	if err != nil {
		logger.Error("[Handler][AppleCallbackV2Handler()]->Error DeleteOauthState()", logger.Field("error", err.Error()))
		c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthLoginRedirectUrl(utilsRedirect.ProviderApple, false, ""))
		return
	}

	appleOAuthUserInfo, err := controller.GetAppleOAuthUserInfo(request.Code, 2)
	if err != nil {
		logger.Error("[Handler][AppleCallbackV2Handler()]->Error GetLineOAuthUserInfo()", logger.Field("error", err.Error()))
		c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthLoginRedirectUrl(utilsRedirect.ProviderApple, false, ""))
		return
	}

	checkAccountOAuth, err := controllerInstance.CheckAccountOAuth(modelDatabase.AccountOAuthProviderApple, modelController.ParamCheckAccountOAuth{
		UserID: appleOAuthUserInfo.UserID,
	})
	if err != nil {
		logger.Error("[Handler][AppleCallbackV2Handler()]->Error CheckAccountOAuth()", logger.Field("error", err.Error()))
		c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthLoginRedirectUrl(utilsRedirect.ProviderApple, false, ""))
		return
	}

	var accountUUID uuid.UUID = checkAccountOAuth.AccountUUID

	if checkAccountOAuth.IsNotFound {
		var userData modelHandler.RequestAppleCallbackUser
		if request.User != "" {
			err := json.Unmarshal([]byte(request.User), &userData)
			if err != nil {
				logger.Error("[Handler][AppleCallbackV2Handler()]->Json unmarshal user data fail", logger.Field("error", err))
				c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthLoginRedirectUrl(utilsRedirect.ProviderApple, false, ""))
				return
			}
		}
		appleOAuthUserInfo.Name = userData.Name.FirstName + " " + userData.Name.LastName
		if appleOAuthUserInfo.Name == " " {
			appleOAuthUserInfo.Name = "User " + appleOAuthUserInfo.Email
		}

		password, err := password.Generate(64, 10, 10, false, false)
		if err != nil {
			logger.Error("[Handler][AppleCallbackV2Handler()]->Error Generate Password", logger.Field("error", err.Error()))
			c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthLoginRedirectUrl(utilsRedirect.ProviderApple, false, ""))
			return
		}
		createAccount, err := controllerInstance.CreateAccount(modelController.ParamCreateAccount{
			Name:               appleOAuthUserInfo.Name,
			Email:              appleOAuthUserInfo.Email,
			Password:           password,
			EmailVerifyApprove: true,
		})
		if err != nil {
			logger.Error("[Handler][AppleCallbackV2Handler()]->Error Create Account", logger.Field("error", err.Error()))
			c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthLoginRedirectUrl(utilsRedirect.ProviderApple, false, ""))
			return
		}

		accountUUID = createAccount.UUID

		err = controllerInstance.CreateAccountOAuth(modelController.ParamCreateAccountOAuth{
			AccountUUID: accountUUID.String(),
			Provider:    modelDatabase.AccountOAuthProviderApple,
			UserID:      appleOAuthUserInfo.UserID,
			UserName:    appleOAuthUserInfo.Name,
			UserEmail:   appleOAuthUserInfo.Email,
			UserPicture: "",
		})
		if err != nil {
			logger.Error("[Handler][AppleCallbackV2Handler()]->Error CreateAccountOAuth()", logger.Field("error", err.Error()))
			c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthLoginRedirectUrl(utilsRedirect.ProviderApple, false, ""))
			return
		}
	}

	createTemporaryCode, err := controllerInstance.CreateTemporaryCode(modelController.ParamTemporaryCode{
		AccountUUID: accountUUID.String(),
		Type:        modelDatabase.TemporaryCodeTypeAuthTokenCode,
	})
	if err != nil {
		logger.Error("[Handler][AppleCallbackV2Handler()]->Error CreateTemporaryCode()", logger.Field("error", err.Error()))
		c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthLoginRedirectUrl(utilsRedirect.ProviderApple, false, ""))
		return
	}

	codeUUIDEncryptBase64, err := controller.EncryptTemporaryCode(createTemporaryCode.CodeUUID.String())
	if err != nil {
		logger.Error("[Handler][AppleCallbackV2Handler()]->Error EncryptTemporaryCode()", logger.Field("error", err.Error()))
		c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthLoginRedirectUrl(utilsRedirect.ProviderApple, false, ""))
		return
	}

	databaseTx.Commit()

	c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthLoginRedirectUrl(utilsRedirect.ProviderApple, true, codeUUIDEncryptBase64))
}

// Connect Apple OAuth
func RequestConnectAppleOAuthHandler(c *gin.Context) {
	var queryParam modelHandler.QueryParamOAuthConnect

	if err := c.ShouldBindQuery(&queryParam); err != nil {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusBadRequest,
		})
		return
	}

	isValidatePass, _, validatorError := PTGUvalidator.Validate(queryParam)
	if validatorError != nil {
		logger.Error("[Handler][RequestConnectAppleOAuthHandler()]->Error Validate Data", logger.Field("error", validatorError.Error()))
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

	var response modelHandler.ResponseRequestConnectOAuth

	databaseTx := database.DB.Begin()
	controllerInstance := controller.NewController(databaseTx)
	defer databaseTx.Rollback()

	checkAccountOAuth, err := controllerInstance.CheckAccountOAuth(modelDatabase.AccountOAuthProviderApple, modelController.ParamCheckAccountOAuth{
		AccountUUID: c.GetString("ACCOUNT_UUID"),
	})
	if err != nil {
		logger.Error("[Handler][RequestConnectAppleOAuthHandler()]->Error CheckAccountOAuth()", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}
	if !checkAccountOAuth.IsNotFound {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusBadRequest,
			Error:        "Apple OAuth Already Connected",
		})
		return
	}

	err = controllerInstance.DeleteTemporaryCode(modelController.ParamTemporaryCode{
		AccountUUID: c.GetString("ACCOUNT_UUID"),
		Type:        modelDatabase.TemporaryCodeTypeOAuthStateApple,
	})
	if err != nil {
		logger.Error("[Handler][RequestConnectAppleOAuthHandler()]->Error DeleteTemporaryCode()", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}
	err = controllerInstance.DeleteTemporaryCode(modelController.ParamTemporaryCode{
		AccountUUID: c.GetString("ACCOUNT_UUID"),
		Type:        modelDatabase.TemporaryCodeTypeAppOAuthStateApple,
	})
	if err != nil {
		logger.Error("[Handler][RequestConnectAppleOAuthHandler()]->Error DeleteTemporaryCode()", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}

	var typeConnect string
	if queryParam.IsApplication {
		typeConnect = modelDatabase.TemporaryCodeTypeAppOAuthStateApple
	} else {
		typeConnect = modelDatabase.TemporaryCodeTypeOAuthStateApple
	}
	createTemporaryCode, err := controllerInstance.CreateTemporaryCode(modelController.ParamTemporaryCode{
		AccountUUID: c.GetString("ACCOUNT_UUID"),
		Type:        typeConnect,
	})
	if err != nil {
		logger.Error("[Handler][RequestConnectAppleOAuthHandler()]->Error CreateTemporaryCode()", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}

	databaseTx.Commit()

	response.AuthURL = controller.GenerateAppleOAuthURL(createTemporaryCode.CodeUUID.String(), 1)

	utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
		ResponseCode: http.StatusOK,
		Data:         response,
	})
}

func RequestDisconnectAppleOAuthHandler(c *gin.Context) {
	databaseTx := database.DB.Begin()
	controllerInstance := controller.NewController(databaseTx)
	defer databaseTx.Rollback()

	deleteOAuthAccount, err := controllerInstance.DeleteAccountOAuth(modelController.ParamDeleteAccountOAuth{
		AccountUUID: c.GetString("ACCOUNT_UUID"),
		Provider:    modelDatabase.AccountOAuthProviderApple,
	})
	if err != nil {
		logger.Error("[Handler][RequestDisconnectAppleOAuthHandler()]->Error DeleteAccountOAuth()", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}

	if deleteOAuthAccount.IsNotFound {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusNotFound,
			Error:        "Apple OAuth Not Connected",
		})
		return
	}

	databaseTx.Commit()

	utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
		ResponseCode: http.StatusOK,
		Data:         "Apple OAuth Disconnected",
	})
}

func AppleCallbackHandler(c *gin.Context) {
	var request modelHandler.RequestAppleCallback

	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	isValidatePass, _, validatorError := PTGUvalidator.Validate(request)
	if validatorError != nil {
		logger.Error("[Handler][AppleCallbackHandler()]->Error Validate Data", logger.Field("error", validatorError.Error()))
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

	var checkType string = modelDatabase.TemporaryCodeTypeOAuthStateApple
	var checkTemporaryCode modelController.ReturnCheckTemporaryCode
	checkTemporaryCode, err := controllerInstance.CheckTemporaryCode(modelController.ParamCheckTemporaryCode{
		CodeUUID: request.State,
		Type:     checkType,
	})
	if err != nil {
		logger.Error("[Handler][AppleCallbackHandler()]->Error CheckTemporaryCode()", logger.Field("error", err.Error()))
		c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthConnectRedirectUrl(utilsRedirect.ProviderApple, false))
		return
	}
	if checkTemporaryCode.IsNotFound {
		checkType = modelDatabase.TemporaryCodeTypeAppOAuthStateApple
		checkTemporaryCode, err = controllerInstance.CheckTemporaryCode(modelController.ParamCheckTemporaryCode{
			CodeUUID: request.State,
			Type:     checkType,
		})
		if err != nil {
			logger.Error("[Handler][AppleCallbackHandler()]->Error CheckTemporaryCode()", logger.Field("error", err.Error()))
			c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthConnectRedirectUrl(utilsRedirect.ProviderApple, false))
			return
		}
	}
	if checkTemporaryCode.IsExpired {
		c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthConnectRedirectUrl(utilsRedirect.ProviderApple, false))
		return
	}

	err = controllerInstance.DeleteTemporaryCode(modelController.ParamTemporaryCode{
		AccountUUID: checkTemporaryCode.AccountUUID.String(),
		Type:        modelDatabase.TemporaryCodeTypeOAuthStateApple,
	})
	if err != nil {
		logger.Error("[Handler][AppleCallbackHandler()]->Error DeleteTemporaryCode()", logger.Field("error", err.Error()))
		c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthConnectRedirectUrl(utilsRedirect.ProviderApple, false))
		return
	}

	appleOAuthUserInfo, err := controller.GetAppleOAuthUserInfo(request.Code, 1)
	if err != nil {
		logger.Error("[Handler][AppleCallbackHandler()]->Error GetLineOAuthUserInfo()", logger.Field("error", err.Error()))
		c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthConnectRedirectUrl(utilsRedirect.ProviderApple, false))
		return
	}
	var userData modelHandler.RequestAppleCallbackUser
	if request.User != "" {
		err := json.Unmarshal([]byte(request.User), &userData)
		if err != nil {
			logger.Error("[Handler][AppleCallbackV2Handler()]->Json unmarshal user data fail", logger.Field("error", err))
			c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthLoginRedirectUrl(utilsRedirect.ProviderApple, false, ""))
			return
		}
	}
	appleOAuthUserInfo.Name = userData.Name.FirstName + " " + userData.Name.LastName
	if appleOAuthUserInfo.Name == " " {
		appleOAuthUserInfo.Name = "User " + appleOAuthUserInfo.Email
	}

	checkAccountOAuth, err := controllerInstance.CheckAccountOAuth(modelDatabase.AccountOAuthProviderApple, modelController.ParamCheckAccountOAuth{
		UserID: appleOAuthUserInfo.UserID,
	})
	if err != nil {
		logger.Error("[Handler][AppleCallbackHandler()]->Error CheckAccountOAuth()", logger.Field("error", err.Error()))
		c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthLoginRedirectUrl(utilsRedirect.ProviderApple, false, ""))
		return
	}

	if !checkAccountOAuth.IsNotFound {
		if checkType == modelDatabase.TemporaryCodeTypeAppOAuthStateApple {
			c.Redirect(http.StatusFound, utilsRedirect.GenerateAppOAuthConnectRedirectUrl(utilsRedirect.ProviderApple, false))
		} else {
			c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthConnectRedirectUrl(utilsRedirect.ProviderApple, false))
		}
		return
	}

	err = controllerInstance.CreateAccountOAuth(modelController.ParamCreateAccountOAuth{
		AccountUUID: checkTemporaryCode.AccountUUID.String(),
		Provider:    modelDatabase.AccountOAuthProviderApple,
		UserID:      appleOAuthUserInfo.UserID,
		UserName:    appleOAuthUserInfo.Name,
		UserEmail:   appleOAuthUserInfo.Email,
		UserPicture: "",
	})
	if err != nil {
		logger.Error("[Handler][AppleCallbackHandler()]->Error CreateAccountOAuth()", logger.Field("error", err.Error()))
		c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthConnectRedirectUrl(utilsRedirect.ProviderApple, false))
		return
	}

	databaseTx.Commit()

	if checkType == modelDatabase.TemporaryCodeTypeAppOAuthStateApple {
		c.Redirect(http.StatusFound, utilsRedirect.GenerateAppOAuthConnectRedirectUrl(utilsRedirect.ProviderApple, true))
	} else {
		c.Redirect(http.StatusFound, utilsRedirect.GenerateOAuthConnectRedirectUrl(utilsRedirect.ProviderApple, true))
	}
}
