package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/parinyapt/prinflix_backend/controller"
	"github.com/parinyapt/prinflix_backend/database"
	"github.com/parinyapt/prinflix_backend/logger"
	modelController "github.com/parinyapt/prinflix_backend/model/controller"
	modelDatabase "github.com/parinyapt/prinflix_backend/model/database"
	modelUtils "github.com/parinyapt/prinflix_backend/model/utils"
	utilsResponse "github.com/parinyapt/prinflix_backend/utils/response"
)

type MiddlewareReceiverArgument struct {
	Admin   bool
}

func NewMiddleware(receiver MiddlewareReceiverArgument) *MiddlewareReceiverArgument {
	return &receiver
}

func (receiver *MiddlewareReceiverArgument) CheckAccount(c *gin.Context) {
	accountUUID := c.GetString("ACCOUNT_UUID")

	controllerInstance := controller.NewController(database.DB)

	accountInfo, err := controllerInstance.GetAccountInfo(modelController.ParamGetAccountInfo{
		AccountUUID: accountUUID,
	})
	if err != nil {
		logger.Error("[Middleware][CheckAccount()]->Error GetAccountInfo()", logger.Field("error", err.Error()))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusInternalServerError,
		})
		c.Abort()
		return
	}

	if accountInfo.IsNotFound {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusForbidden,
			ErrorCode:    "ACMW0010",
			Error:        "Account Access Denied",
		})
		c.Abort()
		return
	}

	if accountInfo.Status != modelDatabase.AccountStatusActive {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusForbidden,
			ErrorCode:    "ACMW0011",
			Error:        "Account Access Denied",
		})
		c.Abort()
		return
	}

	if !accountInfo.EmailVerified {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusForbidden,
			ErrorCode:    "ACMW0012",
			Error:        "Account Access Denied",
		})
		c.Abort()
		return
	}

	if receiver.Admin {
		if accountInfo.Role != modelDatabase.AccountRoleAdmin {
			utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
				ResponseCode: http.StatusForbidden,
				ErrorCode:    "ACMW0021",
				Error:        "Account Access Denied",
			})
			c.Abort()
			return
		}
	}

	c.Next()
}