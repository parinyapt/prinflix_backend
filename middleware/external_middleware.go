package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	PTGUvalidator "github.com/parinyapt/golang_utils/validator/v1"

	"github.com/parinyapt/prinflix_backend/controller"
	"github.com/parinyapt/prinflix_backend/database"
	"github.com/parinyapt/prinflix_backend/logger"
	modelUtils "github.com/parinyapt/prinflix_backend/model/utils"
	utilsResponse "github.com/parinyapt/prinflix_backend/utils/response"
)

type HeaderAuthorizationToken struct {
	Authorization string `header:"Authorization" validate:"required,startswith=Bearer "`
}

func GetHeaderAuthorizationToken(c *gin.Context) {
	var header HeaderAuthorizationToken

	if err := c.ShouldBindHeader(&header); err != nil {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusBadRequest,
		})
		c.Abort()
		return
	}

	isValidatePass, _, validatorError := PTGUvalidator.Validate(header)
	if validatorError != nil {
		logger.Error("[Middleware][GetHeaderAuthorizationToken()]->Error Validate Data", logger.Field("error", validatorError.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusBadRequest,
		})
		c.Abort()
		return
	}
	if !isValidatePass {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusBadRequest,
		})
		c.Abort()
		return
	}

	token := strings.Split(header.Authorization, " ")[1]

	c.Set("AUTHORIZATION_TOKEN", token)
	c.Next()
}

type AuthWithToken struct {
	Token string `json:"token" validate:"required,jwt"`
}

func AuthWithAccessToken(c *gin.Context) {
	var auth AuthWithToken
	auth.Token = c.GetString("AUTHORIZATION_TOKEN")

	isValidatePass, _, validatorError := PTGUvalidator.Validate(auth)
	if validatorError != nil {
		logger.Error("[Middleware][AuthWithAccessToken()]->Error Validate Data", logger.Field("error", validatorError.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusBadRequest,
		})
		c.Abort()
		return
	}
	if !isValidatePass {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusBadRequest,
		})
		c.Abort()
		return
	}

	tokenInfo, err := controller.ValidateAccessToken(auth.Token)
	if err != nil {
		logger.Error("[Middleware][AuthWithAccessToken()]->Error Validate Access Token", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusUnauthorized,
		})
		c.Abort()
		return
	}

	if tokenInfo.IsExpired {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusUnauthorized,
		})
		c.Abort()
		return
	}

	controllerInstance := controller.NewController(database.DB)

	checkAuthSession, err := controllerInstance.CheckAuthSession(tokenInfo.SessionUUID)
	if err != nil {
		logger.Error("[Middleware][AuthWithAccessToken()]->Error Check Auth Session", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		c.Abort()
		return
	}

	if checkAuthSession.IsExpired || checkAuthSession.IsNotFound {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusUnauthorized,
		})
		c.Abort()
		return
	}

	c.Set("ACCOUNT_UUID", checkAuthSession.AccountUUID.String())
	c.Set("SESSION_UUID", checkAuthSession.SessionUUID.String())
	c.Next()
}
