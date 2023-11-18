package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	PTGUvalidator "github.com/parinyapt/golang_utils/validator/v1"

	"github.com/parinyapt/prinflix_backend/logger"
	modelUtils "github.com/parinyapt/prinflix_backend/model/utils"
	utilsResponse "github.com/parinyapt/prinflix_backend/utils/response"
)

type HeaderAPIKey struct {
	Authorization string `header:"Authorization" validate:"required"`
}

func GetHeaderAPIKey(c *gin.Context) {
	var header HeaderAPIKey

	if err := c.ShouldBindHeader(&header); err != nil {
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusBadRequest,
		})
		c.Abort()
		return
	}

	isValidatePass, _, validatorError := PTGUvalidator.Validate(header)
	if validatorError != nil {
		logger.Error("[Middleware][GetHeaderAPIKey()]->Error Validate Data", logger.Field("error", validatorError.Error()))
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

	c.Set("API_KEY", header.Authorization)
	c.Next()
}

func AuthWithAPIKey(c *gin.Context) {
	apiKey := c.GetString("API_KEY")

	if apiKey != os.Getenv("INTERNAL_API_KEY") {
		logger.Error("[Middleware][AuthWithAPIKey()]->Error API Key", logger.Field("error", "API Key is not valid"))
		utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
			ResponseCode: http.StatusUnauthorized,
		})
		c.Abort()
		return
	}

	c.Next()
}
