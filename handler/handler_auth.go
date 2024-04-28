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
	modelDatabase "github.com/parinyapt/prinflix_backend/model/database"
	modelHandler "github.com/parinyapt/prinflix_backend/model/handler"
	modelUtils "github.com/parinyapt/prinflix_backend/model/utils"
	utilsResponse "github.com/parinyapt/prinflix_backend/utils/response"
)

func LoginHandler(c *gin.Context) {
	var request modelHandler.RequestLogin
	var response modelHandler.ResponseAccessToken

	if err := c.ShouldBindJSON(&request); err != nil {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusBadRequest,
		})
		return
	}

	isValidatePass, errorFieldList, validatorError := PTGUvalidator.Validate(request)
	if validatorError != nil {
		logger.Error("[Handler][LoginHandler()]->Error Validate Data", logger.Field("error", validatorError.Error()))
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

	checkLogin, err := controllerInstance.CheckLogin(modelController.ParamCheckLogin{
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		logger.Error("[Handler][LoginHandler()]->Error Check Login", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}
	if checkLogin.IsNotFound || checkLogin.IsPasswordNotMatch {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusNotFound,
			Error:        "Invalid Email or Password",
		})
		return
	}

	clearAuthSessionErr := controllerInstance.DeleteAuthSessionByAccountUUID(checkLogin.AccountUUID.String())
	if clearAuthSessionErr != nil {
		logger.Error("[Handler][LoginHandler()]->Error Clear Auth Session", logger.Field("error", clearAuthSessionErr.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}

	createAuthSession, err := controllerInstance.CreateAuthSession(checkLogin.AccountUUID.String())
	if err != nil {
		logger.Error("[Handler][LoginHandler()]->Error Create Auth Session", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}

	generateAccessToken, err := controller.GenerateAccessToken(modelController.ParamGenerateAccessToken{
		SessionUUID: createAuthSession.SessionUUID.String(),
		ExpiredAt:   createAuthSession.ExpiredAt,
	})
	if err != nil {
		logger.Error("[Handler][LoginHandler()]->Error Generate Access Token", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}

	databaseTx.Commit()

	response.TokenType = generateAccessToken.TokenType
	response.AccessToken = generateAccessToken.AccessToken
	response.AccessTokenExpireIn = time.Duration(createAuthSession.ExtiredIn.Seconds())

	utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
		ResponseCode: http.StatusOK,
		Data:         response,
	})
}

func RegisterHandler(c *gin.Context) {
	var request modelHandler.RequestRegister

	if err := c.ShouldBindJSON(&request); err != nil {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusBadRequest,
		})
		return
	}

	isValidatePass, errorFieldList, validatorError := PTGUvalidator.Validate(request)
	if validatorError != nil {
		logger.Error("[Handler][RegisterHandler()]->Error Validate Data", logger.Field("error", validatorError.Error()))
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

	createAccount, err := controllerInstance.CreateAccount(modelController.ParamCreateAccount{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		logger.Error("[Handler][RegisterHandler()]->Error Create Account", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}
	if createAccount.IsExist {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusConflict,
			Error:        "Email already exist",
		})
		return
	}

	createTemporaryCode, err := controllerInstance.CreateTemporaryCode(modelController.ParamTemporaryCode{
		AccountUUID: createAccount.UUID.String(),
		Type:        modelDatabase.TemporaryCodeTypeEmailVerification,
	})
	if err != nil {
		logger.Error("[Handler][RegisterHandler()]->Error Create Temporary Code", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}

	codeUUIDEncryptBase64, err := controller.EncryptTemporaryCode(createTemporaryCode.CodeUUID.String())
	if err != nil {
		logger.Error("[Handler][RegisterHandler()]->Error Encrypt Temporary Code", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}

	err = controller.SendEmail(modelController.ParamSendEmail{
		Email: request.Email,
		Data:  codeUUIDEncryptBase64,
		Type:  modelDatabase.TemporaryCodeTypeEmailVerification,
	})
	if err != nil {
		logger.Error("[Handler][RegisterHandler()]->Error Send Email", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}

	databaseTx.Commit()

	utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
		ResponseCode: http.StatusOK,
		Data:         "Register Success",
	})
}

func VerifyTokenHandler(c *gin.Context) {
	var response modelHandler.ResponseVerifyToken

	controllerInstance := controller.NewController(database.DB)

	accountInfo, err := controllerInstance.GetAccountInfo(modelController.ParamGetAccountInfo{
		AccountUUID: c.GetString("ACCOUNT_UUID"),
	})
	if err != nil {
		logger.Error("[Handler][VerifyTokenHandler()]->Error GetAccountInfo()", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}

	if accountInfo.IsNotFound {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}

	googleAccountOauth, err := controllerInstance.CheckAccountOAuth(modelDatabase.AccountOAuthProviderGoogle, modelController.ParamCheckAccountOAuth{
		AccountUUID: c.GetString("ACCOUNT_UUID"),
	})
	if err != nil {
		logger.Error("[Handler][VerifyTokenHandler()]->Error CheckAccountOAuth() - Google", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}
	if !googleAccountOauth.IsNotFound {
		response.OAuth.Google.IsConnected = true
		response.OAuth.Google.Name = googleAccountOauth.Name
		response.OAuth.Google.Picture = googleAccountOauth.Picture
	}

	lineAccountOauth, err := controllerInstance.CheckAccountOAuth(modelDatabase.AccountOAuthProviderLine, modelController.ParamCheckAccountOAuth{
		AccountUUID: c.GetString("ACCOUNT_UUID"),
	})
	if err != nil {
		logger.Error("[Handler][VerifyTokenHandler()]->Error CheckAccountOAuth() - Line", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}
	if !lineAccountOauth.IsNotFound {
		response.OAuth.Line.IsConnected = true
		response.OAuth.Line.Name = lineAccountOauth.Name
		response.OAuth.Line.Picture = lineAccountOauth.Picture
	}

	appleAccountOauth, err := controllerInstance.CheckAccountOAuth(modelDatabase.AccountOAuthProviderApple, modelController.ParamCheckAccountOAuth{
		AccountUUID: c.GetString("ACCOUNT_UUID"),
	})
	if err != nil {
		logger.Error("[Handler][VerifyTokenHandler()]->Error CheckAccountOAuth() - Apple", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}
	if !appleAccountOauth.IsNotFound {
		response.OAuth.Apple.IsConnected = true
		response.OAuth.Apple.Name = appleAccountOauth.Name
		response.OAuth.Apple.Picture = ""
	}

	response.Name = accountInfo.Name
	response.Email = accountInfo.Email
	response.EmailVerified = accountInfo.EmailVerified
	response.Status = accountInfo.Status
	response.ImageStatus = accountInfo.Image
	response.Role = accountInfo.Role
	response.SessionUUID = c.GetString("SESSION_UUID")

	utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
		ResponseCode: http.StatusOK,
		Data:         response,
	})
}

func RefreshTokenHandler(c *gin.Context) {
	var response modelHandler.ResponseAccessToken

	databaseTx := database.DB.Begin()
	controllerInstance := controller.NewController(databaseTx)
	defer databaseTx.Rollback()

	clearAuthSessionErr := controllerInstance.DeleteAuthSessionByAccountUUID(c.GetString("ACCOUNT_UUID"))
	if clearAuthSessionErr != nil {
		logger.Error("[Handler][RefreshTokenHandler()]->Error Clear Auth Session", logger.Field("error", clearAuthSessionErr.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}

	createAuthSession, err := controllerInstance.CreateAuthSession(c.GetString("ACCOUNT_UUID"))
	if err != nil {
		logger.Error("[Handler][RefreshTokenHandler()]->Error Create Auth Session", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}

	generateAccessToken, err := controller.GenerateAccessToken(modelController.ParamGenerateAccessToken{
		SessionUUID: createAuthSession.SessionUUID.String(),
		ExpiredAt:   createAuthSession.ExpiredAt,
	})
	if err != nil {
		logger.Error("[Handler][RefreshTokenHandler()]->Error Generate Access Token", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}

	databaseTx.Commit()

	response.TokenType = generateAccessToken.TokenType
	response.AccessToken = generateAccessToken.AccessToken
	response.AccessTokenExpireIn = time.Duration(createAuthSession.ExtiredIn.Seconds())

	utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
		ResponseCode: http.StatusOK,
		Data:         response,
	})
}

func RevokeTokenHandler(c *gin.Context) {
	databaseTx := database.DB.Begin()
	controllerInstance := controller.NewController(databaseTx)
	defer databaseTx.Rollback()

	err := controllerInstance.DeleteAuthSessionByAccountUUID(c.GetString("ACCOUNT_UUID"))
	if err != nil {
		logger.Error("[Handler][RevokeTokenHandler()]->Error Delete Auth Session", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}

	databaseTx.Commit()

	utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
		ResponseCode: http.StatusOK,
		Data:         "Revoke Token Success",
	})
}

func InternalGoogleLoginHandler(c *gin.Context) {
	var request modelHandler.RequestInternalOAuthLogin
	var response modelHandler.ResponseAccessToken

	if err := c.ShouldBindJSON(&request); err != nil {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusBadRequest,
		})
		return
	}

	isValidatePass, errorFieldList, validatorError := PTGUvalidator.Validate(request)
	if validatorError != nil {
		logger.Error("[Handler][InternalGoogleLoginHandler()]->Error Validate Data", logger.Field("error", validatorError.Error()))
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

	checkAccountOAuth, err := controllerInstance.CheckAccountOAuth(modelDatabase.AccountOAuthProviderGoogle, modelController.ParamCheckAccountOAuth{
		UserID: request.UserID,
	})
	if err != nil {
		logger.Error("[Handler][InternalGoogleLoginHandler()]->Error CheckAccountOAuth()", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}
	if checkAccountOAuth.IsNotFound {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusNotFound,
			Error:        "Account OAuth Not Found or Not Connected",
		})
		return
	}

	clearAuthSessionErr := controllerInstance.DeleteAuthSessionByAccountUUID(checkAccountOAuth.AccountUUID.String())
	if clearAuthSessionErr != nil {
		logger.Error("[Handler][InternalGoogleLoginHandler()]->Error Clear Auth Session", logger.Field("error", clearAuthSessionErr.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}

	createAuthSession, err := controllerInstance.CreateAuthSession(checkAccountOAuth.AccountUUID.String())
	if err != nil {
		logger.Error("[Handler][InternalGoogleLoginHandler()]->Error Create Auth Session", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}

	generateAccessToken, err := controller.GenerateAccessToken(modelController.ParamGenerateAccessToken{
		SessionUUID: createAuthSession.SessionUUID.String(),
		ExpiredAt:   createAuthSession.ExpiredAt,
	})
	if err != nil {
		logger.Error("[Handler][InternalGoogleLoginHandler()]->Error Generate Access Token", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}

	databaseTx.Commit()

	response.TokenType = generateAccessToken.TokenType
	response.AccessToken = generateAccessToken.AccessToken
	response.AccessTokenExpireIn = time.Duration(createAuthSession.ExtiredIn.Seconds())

	utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
		ResponseCode: http.StatusOK,
		Data:         response,
	})
}

func InternalLineLoginHandler(c *gin.Context) {
	var request modelHandler.RequestInternalOAuthLogin
	var response modelHandler.ResponseAccessToken

	if err := c.ShouldBindJSON(&request); err != nil {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusBadRequest,
		})
		return
	}

	isValidatePass, errorFieldList, validatorError := PTGUvalidator.Validate(request)
	if validatorError != nil {
		logger.Error("[Handler][InternalLineLoginHandler()]->Error Validate Data", logger.Field("error", validatorError.Error()))
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

	checkAccountOAuth, err := controllerInstance.CheckAccountOAuth(modelDatabase.AccountOAuthProviderLine, modelController.ParamCheckAccountOAuth{
		UserID: request.UserID,
	})
	if err != nil {
		logger.Error("[Handler][InternalLineLoginHandler()]->Error CheckAccountOAuth()", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}
	if checkAccountOAuth.IsNotFound {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusNotFound,
			Error:        "Account OAuth Not Found or Not Connected",
		})
		return
	}

	clearAuthSessionErr := controllerInstance.DeleteAuthSessionByAccountUUID(checkAccountOAuth.AccountUUID.String())
	if clearAuthSessionErr != nil {
		logger.Error("[Handler][InternalLineLoginHandler()]->Error Clear Auth Session", logger.Field("error", clearAuthSessionErr.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}

	createAuthSession, err := controllerInstance.CreateAuthSession(checkAccountOAuth.AccountUUID.String())
	if err != nil {
		logger.Error("[Handler][InternalLineLoginHandler()]->Error Create Auth Session", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}

	generateAccessToken, err := controller.GenerateAccessToken(modelController.ParamGenerateAccessToken{
		SessionUUID: createAuthSession.SessionUUID.String(),
		ExpiredAt:   createAuthSession.ExpiredAt,
	})
	if err != nil {
		logger.Error("[Handler][InternalLineLoginHandler()]->Error Generate Access Token", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}

	databaseTx.Commit()

	response.TokenType = generateAccessToken.TokenType
	response.AccessToken = generateAccessToken.AccessToken
	response.AccessTokenExpireIn = time.Duration(createAuthSession.ExtiredIn.Seconds())

	utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
		ResponseCode: http.StatusOK,
		Data:         response,
	})
}

func ExchangeCodeToAuthTokenV2Handler(c *gin.Context) {
	var request modelHandler.RequestExchangeCodeToAuthToken
	var response modelHandler.ResponseAccessToken

	if err := c.ShouldBindJSON(&request); err != nil {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusBadRequest,
		})
		return
	}

	isValidatePass, errorFieldList, validatorError := PTGUvalidator.Validate(request)
	if validatorError != nil {
		logger.Error("[Handler][ExchangeCodeToAuthTokenV2Handler()]->Error Validate Data", logger.Field("error", validatorError.Error()))
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

	codeUUID, err := controller.DecryptTemporaryCode(request.Code)
	if err != nil {
		logger.Error("[Handler][ExchangeCodeToAuthTokenV2Handler()]->Error DecryptTemporaryCode()", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusBadRequest,
		})
		return
	}

	checkTemporaryCode, err := controllerInstance.CheckTemporaryCode(modelController.ParamCheckTemporaryCode{
		CodeUUID: codeUUID,
		Type:     modelDatabase.TemporaryCodeTypeAuthTokenCode,
	})
	if err != nil {
		logger.Error("[Handler][ExchangeCodeToAuthTokenV2Handler()]->Error CheckTemporaryCode()", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}
	if checkTemporaryCode.IsNotFound || checkTemporaryCode.IsExpired {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusBadRequest,
			Error:        "Invalid Code or Code Expired",
		})
		return
	}

	err = controllerInstance.DeleteTemporaryCode(modelController.ParamTemporaryCode{
		AccountUUID: checkTemporaryCode.AccountUUID.String(),
		Type:        modelDatabase.TemporaryCodeTypeAuthTokenCode,
	})
	if err != nil {
		logger.Error("[Handler][ExchangeCodeToAuthTokenV2Handler()]->Error DeleteTemporaryCode()", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}

	clearAuthSessionErr := controllerInstance.DeleteAuthSessionByAccountUUID(checkTemporaryCode.AccountUUID.String())
	if clearAuthSessionErr != nil {
		logger.Error("[Handler][ExchangeCodeToAuthTokenV2Handler()]->Error Clear Auth Session", logger.Field("error", clearAuthSessionErr.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}

	createAuthSession, err := controllerInstance.CreateAuthSession(checkTemporaryCode.AccountUUID.String())
	if err != nil {
		logger.Error("[Handler][ExchangeCodeToAuthTokenV2Handler()]->Error Create Auth Session", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}

	generateAccessToken, err := controller.GenerateAccessToken(modelController.ParamGenerateAccessToken{
		SessionUUID: createAuthSession.SessionUUID.String(),
		ExpiredAt:   createAuthSession.ExpiredAt,
	})
	if err != nil {
		logger.Error("[Handler][ExchangeCodeToAuthTokenV2Handler()]->Error Generate Access Token", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}

	databaseTx.Commit()

	response.TokenType = generateAccessToken.TokenType
	response.AccessToken = generateAccessToken.AccessToken
	response.AccessTokenExpireIn = time.Duration(createAuthSession.ExtiredIn.Seconds())

	utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
		ResponseCode: http.StatusOK,
		Data:         response,
	})
}