package handler

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	modelUtils "github.com/parinyapt/prinflix_backend/model/utils"
	utilsResponse "github.com/parinyapt/prinflix_backend/utils/response"
)

func NoRouteHandler(c *gin.Context) {
	utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
		ResponseCode: http.StatusNotFound,
	})
}

type HealthCheckHandlerResponse struct {
	Version string `json:"version"`
}

func HealthCheckHandler(c *gin.Context) {
	utilsResponse.ApiResponse(c, modelUtils.ApiResponseStruct{
		ResponseCode: http.StatusOK,
		Data: HealthCheckHandlerResponse{
			Version: os.Getenv("DEPLOY_VERSION"),
		},
	})
}