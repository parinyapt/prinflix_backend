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
	utilsConfigFile "github.com/parinyapt/prinflix_backend/utils/config_file"
	utilsResponse "github.com/parinyapt/prinflix_backend/utils/response"
)

func EmailVerifyHandler(c *gin.Context) {
	var uriParam modelHandler.UriParamEmailVerifyHandler

	if err := c.ShouldBindUri(&uriParam); err != nil {
		c.Redirect(http.StatusFound, utilsConfigFile.GetFrontendBaseURL()+utilsConfigFile.GetRedirectPagePath(utilsConfigFile.EmailVerifyFailPagePath))
		return
	}

	isValidatePass, _, validatorError := PTGUvalidator.Validate(uriParam)
	if validatorError != nil {
		logger.Error("[Handler][EmailVerifyHandler()]->Error Validate Data", logger.Field("error", validatorError.Error()))
		c.Redirect(http.StatusFound, utilsConfigFile.GetFrontendBaseURL()+utilsConfigFile.GetRedirectPagePath(utilsConfigFile.EmailVerifyFailPagePath))
		return
	}
	if !isValidatePass {
		c.Redirect(http.StatusFound, utilsConfigFile.GetFrontendBaseURL()+utilsConfigFile.GetRedirectPagePath(utilsConfigFile.EmailVerifyFailPagePath))
		return
	}

	databaseTx := database.DB.Begin()
	controllerInstance := controller.NewController(databaseTx)
	defer databaseTx.Rollback()

	codeUUID, err := controller.DecryptTemporaryCode(uriParam.Code)
	if err != nil {
		logger.Error("[Handler][EmailVerifyHandler()]->Error DecryptTemporaryCode()", logger.Field("error", err.Error()))
		c.Redirect(http.StatusFound, utilsConfigFile.GetFrontendBaseURL()+utilsConfigFile.GetRedirectPagePath(utilsConfigFile.EmailVerifyFailPagePath))
		return
	}

	checkTemporaryCode, err := controllerInstance.CheckTemporaryCode(modelController.ParamCheckTemporaryCode{
		CodeUUID: codeUUID,
		Type:     modelDatabase.TemporaryCodeTypeEmailVerification,
	})
	if err != nil {
		logger.Error("[Handler][EmailVerifyHandler()]->Error CheckTemporaryCode()", logger.Field("error", err.Error()))
		c.Redirect(http.StatusFound, utilsConfigFile.GetFrontendBaseURL()+utilsConfigFile.GetRedirectPagePath(utilsConfigFile.EmailVerifyFailPagePath))
		return
	}
	if checkTemporaryCode.IsNotFound || checkTemporaryCode.IsExpired {
		c.Redirect(http.StatusFound, utilsConfigFile.GetFrontendBaseURL()+utilsConfigFile.GetRedirectPagePath(utilsConfigFile.EmailVerifyFailPagePath))
		return
	}

	err = controllerInstance.DeleteTemporaryCode(modelController.ParamTemporaryCode{
		AccountUUID: checkTemporaryCode.AccountUUID.String(),
		Type:        modelDatabase.TemporaryCodeTypeEmailVerification,
	})
	if err != nil {
		logger.Error("[Handler][EmailVerifyHandler()]->Error DeleteTemporaryCode()", logger.Field("error", err.Error()))
		c.Redirect(http.StatusFound, utilsConfigFile.GetFrontendBaseURL()+utilsConfigFile.GetRedirectPagePath(utilsConfigFile.EmailVerifyFailPagePath))
		return
	}

	updateAccount, err := controllerInstance.UpdateAccount(checkTemporaryCode.AccountUUID.String(), modelController.ParamUpdateAccount{
		EmailVerified: true,
	})
	if err != nil {
		logger.Error("[Handler][EmailVerifyHandler()]->Error UpdateAccount()", logger.Field("error", err.Error()))
		c.Redirect(http.StatusFound, utilsConfigFile.GetFrontendBaseURL()+utilsConfigFile.GetRedirectPagePath(utilsConfigFile.EmailVerifyFailPagePath))
		return
	}
	if updateAccount.IsNotFound {
		c.Redirect(http.StatusFound, utilsConfigFile.GetFrontendBaseURL()+utilsConfigFile.GetRedirectPagePath(utilsConfigFile.EmailVerifyFailPagePath))
		return
	}

	databaseTx.Commit()

	c.Redirect(http.StatusFound, utilsConfigFile.GetFrontendBaseURL()+utilsConfigFile.GetRedirectPagePath(utilsConfigFile.EmailVerifySuccessPagePath))
}

func RequestEmailVerifyHandler(c *gin.Context) {
	databaseTx := database.DB.Begin()
	controllerInstance := controller.NewController(databaseTx)
	defer databaseTx.Rollback()

	accountInfo, err := controllerInstance.GetAccountInfo(modelController.ParamGetAccountInfo{
		AccountUUID: c.GetString("ACCOUNT_UUID"),
	})
	if err != nil {
		logger.Error("[Handler][RequestForgotPasswordHandler()]->Error GetAccountInfo()", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}

	if accountInfo.IsNotFound {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusBadRequest,
		})
		return
	}

	if accountInfo.EmailVerified {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusOK,
			Data:         "Email already verified",
		})
		return
	}

	err = controllerInstance.DeleteTemporaryCode(modelController.ParamTemporaryCode{
		AccountUUID: accountInfo.AccountUUID.String(),
		Type:        modelDatabase.TemporaryCodeTypeEmailVerification,
	})
	if err != nil {
		logger.Error("[Handler][RequestEmailVerifyHandler()]->Error DeleteTemporaryCode()", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}

	createTemporaryCode, err := controllerInstance.CreateTemporaryCode(modelController.ParamTemporaryCode{
		AccountUUID: accountInfo.AccountUUID.String(),
		Type:        modelDatabase.TemporaryCodeTypeEmailVerification,
	})
	if err != nil {
		logger.Error("[Handler][RequestEmailVerifyHandler()]->Error CreateTemporaryCode()", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}

	codeUUIDEncryptBase64, err := controller.EncryptTemporaryCode(createTemporaryCode.CodeUUID.String())
	if err != nil {
		logger.Error("[Handler][RequestEmailVerifyHandler()]->Error EncryptTemporaryCode()", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}

	err = controller.SendEmail(modelController.ParamSendEmail{
		Email: accountInfo.Email,
		Data:  codeUUIDEncryptBase64,
		Type:  modelDatabase.TemporaryCodeTypeEmailVerification,
	})
	if err != nil {
		logger.Error("[Handler][RequestEmailVerifyHandler()]->Error SendEmail()", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		return
	}

	databaseTx.Commit()

	utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
		ResponseCode: http.StatusOK,
		Data:         "Please check your email to verify your email",
	})
}