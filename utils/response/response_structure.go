package utilsResponse

import (
	"time"

	"github.com/gin-gonic/gin"
	
	modelUtils "github.com/parinyapt/prinflix_backend/model/utils"
)

func JsonResponse(c *gin.Context, param modelUtils.JsonResponseStruct) {
	c.JSON(param.ResponseCode, modelUtils.JsonResponseStructDetail{
		Timestamp: time.Now().Format(time.RFC3339),
		Success:   param.Detail.Success,
		ErrorCode: param.Detail.ErrorCode,
		Data:      param.Detail.Data,
		Error:     param.Detail.Error,
	})
}
